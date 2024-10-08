import { lex } from "./Lexer.ts";
import type { MdNode } from "./Markdown.interface.ts";
import { parse } from "./Parser.ts";


export function ConvertToHtmlTree(markdown: string): MdNode {
  const lexed = lex(markdown);
  const parsed = parse(lexed);
  return parsed;
}