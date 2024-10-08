import type { ExtendedPageContent } from "../../types/BasicPageContent.interface.ts";
import { pageContent } from "./content.ts";
import type { MdNode } from "$lib/markdown/Markdown.interface.ts";
import { ConvertToHtmlTree } from "$lib/markdown/Markdown.ts";


/** @type {import('./$types').PageLoad} */
export const load = async () => {
  const parsedPage: ExtendedPageContent[] = pageContent.map((c) => {
    const md = ConvertToHtmlTree(c.MarkdownContent).content as MdNode[];
    console.log(md);
    return {
      ...c,
      ParsedMarkdown: md
    } as ExtendedPageContent
  })

  return {
      data: parsedPage
  };
};
