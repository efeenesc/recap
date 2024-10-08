import MarkdownRenderer from "../../components/markdown-renderer/MarkdownRenderer.svelte";
import type { BasicPageContent } from "../../types/BasicPageContent.interface.ts";

const recap = {
Title: "Recap",
MarkdownContent: `
Recap is a utility for keeping and tracking reports about your daily activity.
Recap works by taking screenshots and sending them to the user's preferred vision API and model for it to describe the user's activity. These
descriptions are later used to generate daily reports upon command, creating a recap of your activity.

It is inspired by Microsoft Recall, which has been the subject of deserved scrutiny due to privacy and security concerns. Recap has a different
objective than Microsoft Recall, yet it has the same privacy and security concerns. It brings these same concerns to other platforms Recap can run on
to the delight of even more privacy-consciented users.

Locally-run vision and text generation models on Ollama should be used with Recap over Gemini. Due to the system demands of local vision and text generation models,
this recommendation can be applied on few devices.

Google asserts that images uploaded to Gemini are not used during its AI models' training, which applies to users of the free Gemini 1.5 Flash model as well.
All images uploaded by Recap are deleted immediately by Recap after a description of the images has been generated. Responses from Gemini containing screenshot
descriptions will be used for training purposes.

By using the Gemini API with Recap, you acknowledge the risk associated with the API.

A more detailed description of Recap can be found on [the project's page](https://github.com/efeenesc/recap).
`,
}

const support = {
Title: "Support",
MarkdownContent: `
For general application support and bug reports, please create an issue on [Recap's GitHub repository](https://github.com/efeenesc/recap).

Please provide details on your operating system and version alongside your issue. For bug reports, please attach screenshots and error messages if possible.
`,
}

const contact = {
Title: "Contact",
MarkdownContent: `
I can be reached from **[hello@efeenescamci.com](mailto:hello@efeenescamci.com)**,
or from my personal **[efeenescamci@icloud.com](mailto:efeenescamci@icloud.com)** e-mail address for any queries about Recap.
`
}

const license = {
Title: "License",
MarkdownContent: `
This project uses version 3 of the GNU Affero General Public License (AGPL).

A copy of the license is available [Recap's GitHub repository](https://github.com/efeenesc/recap/blob/master/LICENSE).

The full license, along with potentially helpful resources, can be found on [GNU's website](https://www.gnu.org/licenses/agpl-3.0.en.html).
`
}

const appVersion = {
MarkdownContent: `
Recap

0.0.1

efeenesc

[efeenescamci.com](efeenescamci.com)
`,
}

const test = {
MarkdownContent: `
blah blah blah
uuuuuuuuugh
`,
}

export const pageContent: BasicPageContent[] = [
    recap,
    support,
    contact,
    license,
    appVersion,
    // test
];

