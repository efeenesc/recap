<div align="center">
<picture>
    <img src="/assets/appicon.png" width="30%">
</picture>
</div>
<h1 align="center">
Recap
</h1>

A take on the idea of Microsoft Recall with the purpose of generating reports of user activity, written in Go.

With its default configuration, it captures a screenshot of all displays every 5 minutes. Every 2 hours, all unprocessed screenshots are sent to the user's preferred vision model to process them and describe the user's activity.

A report can be generated by the user at any time by selecting screenshots through the user interface and clicking the 'Report' button. A report will be generated with screenshots from today, which will be accessible from the 'Reports' page.

## API Support

Recap currently works with Ollama and Gemini APIs. Gemini 1.5 Flash (although risky) is an easy and free way of using this program. Please read the **Warning** section below for details on the security risks.

## Warning

__This is a dangerous idea, vulnerable in varying degrees to the same problems Microsoft Recall has. Gemini 1.5 Flash's free tier, the default for this project, does not have enterprise-level data privacy guarantees.__

All submitted screenshots to the Gemini API are manually deleted after descriptions are generated to reduce risk. [Google clarifies that uploaded files or pixels of uploaded images are not used to train their models unless the user provides feedback, which this project does not.](https://support.google.com/gemini/answer/13594961?hl=en#uploaded_images) Conversations on Gemini 1.5 Flash may be processed and stored by Google - this project allows the use of different models and APIs for vision and text generation to further reduce risk. This means that it is possible to send screenshots to Gemini for its vision model, and use Ollama for text generation.

__Use of a locally-run vision and text model, although perhaps not realistic for all devices, is recommended.__ Ollama API can be used to run local models by changing app settings.

Screenshots captured by the program are saved inside the user-specified folder. The database file is not encrypted.

## Instructions

Change the default configuration in settings to fit your purposes. You can change:
 - The text and vision API connectors and models to be used (only Ollama and Gemini connectors exist at this time)
 - Path where screenshots will be saved
 - Minute interval between screenshots
 - Minute interval between sending screenshots to the vision model
 - Prompts used during screenshot description generation and report generation

If using Gemini, generate your API key from Google's AI Console and set it in the app's settings.

### Limitations and considerations

Only Ollama and Gemini APIs are supported at this time.

The Ollama model you choose has to be of good enough quality to decipher text from the screenshots; reports produced by the program will be highly inaccurate otherwise.

Gemini's free plan is enough for the program's core functionality to work at a satisfying accuracy. Gemini's free 1.5 Flash model allows 1,500 requests to be sent in a day. Setting the screenshot interval too low may cause hitting the limit prematurely. At the default rate of one screenshot every 5 minutes, this program runs its core functionality well below the limit, at 288 requests a day. Setting the interval to 1 minute has the program send 1,440 requests in a day, which is the lowest recommended interval for users of the free plan. If setting screenshot interval to 1 minute, make sure to set the screenshot API schedule interval to 15 minutes or below. Doing so will send screenshots to Gemini every 15 minutes for processing to keep up with the rate at which screenshots are generated.

### Build

Use `wails build`.

## Technical Description

This project was built using Go and the Wails GUI framework. SvelteKit, TypeScript are used for the frontend. Support for other APIs may be added by creating a struct that implements the VisionTextModel interface inside a new file called `internal/models/{modelname}/{modelname}.go`. Change the Initialize function inside `internal/llm/llm.go` to support the new API.

User settings are stored in the SQLite database. Defaults for user settings, what user settings are allowed, and more, are defined inside `internal/db/settings.go`.

Screenshots are saved as PNG with stdlib's image/png library set to best compression. JPEG thumbnails are generated with quality set to 40%.
