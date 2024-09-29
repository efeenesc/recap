export type MdNodeType =
    | "text"
    | "line"
    | "p"
    | "ul"
    | "ol"
    | "li"
    | "h1"
    | "h2"
    | "h3"
    | "h4"
    | "h5"
    | "h6"
    | "code"
    | "bq"
    | "b"
    | "i"
    | "bi"
    | "a"
    | "st"
    | "br"
    | "document";

export interface IDictionary<TValue> {
    [id: string]: TValue;
}

export const MdNodeDict: IDictionary<string> = {
    "#": "h1",
    "##": "h2",
    "###": "h3",
    "####": "h4",
    "#####": "h5",
    "######": "h6",
    "*": "i",
    "**": "b",
    "***": "bi",
    "~~": "st",
    "`": "code",
    "-": "ul",
    ">": "bq",
};

export const typeMap: { [key: string]: MdNodeType } = {
    "#": "h1",
    "##": "h2",
    "###": "h3",
    "####": "h4",
    "#####": "h5",
    "######": "h6",
    "-": "ul",
    ">": "bq",
    "*": "i",
    "**": "b",
    "***": "bi",
    _: "i",
    __: "b",
    ___: "bi",
    "`": "code",
    "~~": "st",
};

export const MdSpecialChar: string[] = [
    "#",
    "*",
    "_",
    "~",
    "-",
    ">",
    "[",
    "]",
    "(",
    ")",
];
export const MdWrapChar: string[] = ["*", "_", "~", "`"];

export class MdNode {
    type: MdNodeType;
    content: MdNode[] | string;
    url?: string;

    constructor(type: MdNodeType, content: MdNode[] | string = []) {
        this.type = type;
        this.content = content;
    }
}

export function lex(mdstr: string): string[][] {
    // Split the input string into lines
    let lines: string[] = mdstr.split("\n");

    // Remove the last empty line if it exists
    if (lines[lines.length - 1] === "") {
        lines.pop();
    }

    // Process each line
    const result = lines.map((line: string) => {
        const tokens: string[] = line.split(" ");
        const processedTokens: string[] = [];

        tokens.forEach((token) => {
            const newTokens: string[] = processToken(token);
            processedTokens.push(...newTokens);
        });

        return processedTokens;
    });

    console.log(result);
    return result;
}

function processToken(token: string): string[] {
    const newTokens: string[] = [];
    let buffer: string[] = [];
    let isSpecial = false;
    let specialRepeat = 0;
    let lastChar: string = "";

    for (let i = 0; i < token.length; i++) {
        const char = token[i];

        if (MdWrapChar.includes(char)) {
            handleSpecialChar(
                char,
                newTokens,
                buffer,
                isSpecial,
                specialRepeat,
                lastChar
            );
            isSpecial = true;
            specialRepeat =
                char === lastChar || specialRepeat === 0
                    ? specialRepeat + 1
                    : 1;
        } else if (isSpecial) {
            handleNonSpecialCharAfterSpecial(
                char,
                newTokens,
                buffer,
                lastChar,
                specialRepeat
            );
            isSpecial = false;
            specialRepeat = 0;
        } else {
            buffer.push(char);
        }

        lastChar = char;
    }

    flushBuffer(newTokens, buffer, isSpecial, lastChar, specialRepeat);

    return newTokens;
}

function handleSpecialChar(
    char: string,
    newTokens: string[],
    buffer: string[],
    isSpecial: boolean,
    specialRepeat: number,
    lastChar: string
) {
    if (isSpecial && char !== lastChar) {
        newTokens.push(lastChar.repeat(specialRepeat));
        flushBuffer(newTokens, buffer);
    }
}

function handleNonSpecialCharAfterSpecial(
    char: string,
    newTokens: string[],
    buffer: string[],
    lastChar: string,
    specialRepeat: number
) {
    flushBuffer(newTokens, buffer);
    newTokens.push(lastChar.repeat(specialRepeat));
    buffer.push(char);
}

function flushBuffer(
    newTokens: string[],
    buffer: string[],
    isSpecial: boolean = false,
    lastChar: string = "",
    specialRepeat: number = 0
) {
    if (buffer.length > 0) {
        newTokens.push(buffer.join(""));
        buffer.length = 0; // Clear the buffer
    }
    if (isSpecial) {
        newTokens.push(lastChar.repeat(specialRepeat));
    }
}

function convertUlToLi(nodes: MdNode[]): MdNode[] {
    return nodes.map((node) =>
        node.type === "ul" ? new MdNode("li", node.content) : node
    );
}

function handleListItems(rootContent: MdNode[], result: { nodes: MdNode[] }) {
    const firstNode = result.nodes[0];
    const listType = firstNode.type === "ul" ? "ul" : "ol";
    const lastNode = rootContent[rootContent.length - 1];

    // Convert 'ul' nodes to 'li' nodes
    const processedNodes = convertUlToLi(result.nodes);

    if (lastNode && lastNode.type === listType) {
        // If the last node is the same type of list, append to it
        (lastNode.content as MdNode[]).push(...processedNodes);
    } else {
        // Otherwise, create a new list node
        const newNode = new MdNode(listType, processedNodes);
        rootContent.push(newNode);
    }
}

export function parse(l: string[][]): MdNode {
    const rootNode: MdNode = new MdNode("document", []);
    let prevIsNewline = false;
    const arrLen = l.length;

    for (let idx = 0; idx < arrLen; idx++) {
        const line = l[idx];
        const rootContent = rootNode.content as MdNode[];

        if (line.length === 0) {
            if (prevIsNewline) {
                prevIsNewline = false;
                rootContent.push(new MdNode("br", ""));
            }

            prevIsNewline = true;
            continue;
        }

        const result = processTokens(line);

        switch (result.nodes[0].type) {
            case "li":
            case "ul":
                handleListItems(rootContent, result);
                break;

            default:
                rootContent.push(...result.nodes);
                break;
        }
    }

    console.log(rootNode);
    return rootNode;
}

function lookAheadFind(
    targetToken: string,
    tokens: string[],
    startFrom: number
): number {
    const arrLen = tokens.length;
    for (; startFrom < arrLen; startFrom++) {
        if (tokens[startFrom] === targetToken) return startFrom;
    }
    return -1;
}

function processTokens(
    tokens: string[],
    index: number = 0,
    closingToken?: string
): { nodes: MdNode[]; index: number } {
    const arrlen = tokens.length;
    const nodes: MdNode[] = [];
    let textBuffer: string[] = [];

    function flushTextBuffer() {
        if (textBuffer.length > 0) {
            nodes.push(new MdNode("text", textBuffer.join(" ")));
            textBuffer = [];
        }
    }

    for (; index < arrlen; index++) {
        const token = tokens[index];

        if (closingToken && token === closingToken) {
            flushTextBuffer();
            return { nodes, index };
        }

        // Get node type. If node type is undefined, then this is not a special Markdown character - likely just text.
        const nodeType = getNodeType(token);
        if (nodeType) {
            // Flush all text content and identify what token this is
            flushTextBuffer();
            const closingToken = getClosingToken(token);
            let tokenClosed = true;
            if (closingToken) {
                tokenClosed =
                    lookAheadFind(closingToken, tokens, index + 1) === -1
                        ? false
                        : true;

                // If the * was not closed, and is the first character in the line, then this is likely just an unordered list.
                // Read child nodes of the list item and add them as list items (<li>), then set the parent element to <ul>
                if (!tokenClosed && index === 0 && ['-', '*', '+'].includes(token)) {
                    const { nodes: childNodes, index: newIndex } =
                        processTokens(tokens, index + 1, closingToken);

                    nodes.push(new MdNode("ul", childNodes));
                    index = newIndex;
                    continue;
                }
            }

            // If token was closed (i.e. * .... *), then process the text content inside tokens
            // This will then assign an element name (i.e. <i>)
            if (tokenClosed) {
                const { nodes: childNodes, index: newIndex } = processTokens(
                    tokens,
                    index + 1,
                    closingToken
                );

                nodes.push(new MdNode(nodeType, childNodes));
                index = newIndex;
                continue;
            }
        }

        textBuffer.push(token);
    }

    flushTextBuffer();
    return { nodes, index };
}

/**
 * Get node type. Match from typeMap. If given numbered list (i.e. 1., 2., etc.), returns <li>
 * @param token
 * @returns
 */
function getNodeType(token: string): MdNodeType | undefined {
    const result = typeMap[token];

    if (result) return result;

    if (
        token.endsWith(".") &&
        !isNaN(Number(token.substring(0, token.length - 1)))
    ) {
        return "li";
    }

    return;
}

/**
 * Get the closing token for a given token if it exists. This is only applicable to wrapped tokens (i.e. ` _ *). Returns undefined if token isn't wrappable
 * @param token
 * @returns
 */
function getClosingToken(token: string): string | undefined {
    switch (token) {
        case "*":
        case "**":
        case "***":
        case "_":
        case "__":
        case "___":
        case "`":
        case "~~":
            return token;

        default:
            return;
    }
}
