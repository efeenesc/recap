import MarkdownRenderer from "../../components/markdown-renderer/MarkdownRenderer.svelte";
import type { BasicPageContent } from "../../types/BasicPageContent.interface.ts";

const recap = {
Title: "Recap",
MarkdownContent: `
Recap is a utility for keeping and tracking reports about your daily activity.
Recap works by taking screenshots and sending them to the user's preferred vision API and model for it to describe the user's activity. These
descriptions are later used to generate daily reports upon command, creating a recap of your activity.

It is recommended to use locally-run vision and text generation models on Ollama with Recap instead of Gemini. Due to the system demands of local vision and text generation models,
this recommendation is applicable to a limited number of devices.

Google asserts that images uploaded to Gemini are not used during its AI models' training, which applies to users of the free Gemini 1.5 Flash model as well.
All images uploaded by Recap are deleted from Gemini immediately by Recap after a description of the images has been generated. However, responses from Gemini containing screenshot
descriptions may be used for training purposes.

By using the Gemini API with Recap, you acknowledge associated risks.

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
I can be reached at **[hello@efeenescamci.com](mailto:hello@efeenescamci.com)**,
or at my personal email address **[efeenescamci@icloud.com](mailto:efeenescamci@icloud.com)** for any queries about Recap.
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

export const pageContent: BasicPageContent[] = [
    recap,
    support,
    contact,
    license,
    appVersion,
];

