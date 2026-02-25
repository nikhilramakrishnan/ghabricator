package diff

// ChromaCSSLight is the Chroma syntax highlighting CSS for the "github" style.
// Generated via: chroma/v2 formatters/html WriteCSS with styles.Get("github").
// These classes are emitted by the formatter configured in highlight.go with
// WithClasses(true). The .chroma prefix matches the class added to the diff table.
const ChromaCSSLight = `/* Chroma — github style (light mode) */
/* Error */ .chroma .err { color: #1f2328 }
/* Keyword */ .chroma .k { color: #cf222e }
/* KeywordConstant */ .chroma .kc { color: #cf222e }
/* KeywordDeclaration */ .chroma .kd { color: #cf222e }
/* KeywordNamespace */ .chroma .kn { color: #cf222e }
/* KeywordPseudo */ .chroma .kp { color: #cf222e }
/* KeywordReserved */ .chroma .kr { color: #cf222e }
/* KeywordType */ .chroma .kt { color: #cf222e }
/* NameAttribute */ .chroma .na { color: #1f2328 }
/* NameClass */ .chroma .nc { color: #1f2328 }
/* NameConstant */ .chroma .no { color: #0550ae }
/* NameDecorator */ .chroma .nd { color: #0550ae }
/* NameEntity */ .chroma .ni { color: #6639ba }
/* NameLabel */ .chroma .nl { color: #990000; font-weight: bold }
/* NameNamespace */ .chroma .nn { color: #24292e }
/* NameOther */ .chroma .nx { color: #1f2328 }
/* NameTag */ .chroma .nt { color: #0550ae }
/* NameBuiltin */ .chroma .nb { color: #6639ba }
/* NameBuiltinPseudo */ .chroma .bp { color: #6a737d }
/* NameVariable */ .chroma .nv { color: #953800 }
/* NameVariableClass */ .chroma .vc { color: #953800 }
/* NameVariableGlobal */ .chroma .vg { color: #953800 }
/* NameVariableInstance */ .chroma .vi { color: #953800 }
/* NameVariableMagic */ .chroma .vm { color: #953800 }
/* NameFunction */ .chroma .nf { color: #6639ba }
/* NameFunctionMagic */ .chroma .fm { color: #6639ba }
/* LiteralString */ .chroma .s { color: #0a3069 }
/* LiteralStringAffix */ .chroma .sa { color: #0a3069 }
/* LiteralStringBacktick */ .chroma .sb { color: #0a3069 }
/* LiteralStringChar */ .chroma .sc { color: #0a3069 }
/* LiteralStringDelimiter */ .chroma .dl { color: #0a3069 }
/* LiteralStringDoc */ .chroma .sd { color: #0a3069 }
/* LiteralStringDouble */ .chroma .s2 { color: #0a3069 }
/* LiteralStringEscape */ .chroma .se { color: #0a3069 }
/* LiteralStringHeredoc */ .chroma .sh { color: #0a3069 }
/* LiteralStringInterpol */ .chroma .si { color: #0a3069 }
/* LiteralStringOther */ .chroma .sx { color: #0a3069 }
/* LiteralStringRegex */ .chroma .sr { color: #0a3069 }
/* LiteralStringSingle */ .chroma .s1 { color: #0a3069 }
/* LiteralStringSymbol */ .chroma .ss { color: #032f62 }
/* LiteralNumber */ .chroma .m { color: #0550ae }
/* LiteralNumberBin */ .chroma .mb { color: #0550ae }
/* LiteralNumberFloat */ .chroma .mf { color: #0550ae }
/* LiteralNumberHex */ .chroma .mh { color: #0550ae }
/* LiteralNumberInteger */ .chroma .mi { color: #0550ae }
/* LiteralNumberIntegerLong */ .chroma .il { color: #0550ae }
/* LiteralNumberOct */ .chroma .mo { color: #0550ae }
/* Operator */ .chroma .o { color: #0550ae }
/* OperatorWord */ .chroma .ow { color: #0550ae }
/* Punctuation */ .chroma .p { color: #1f2328 }
/* Comment */ .chroma .c { color: #57606a }
/* CommentHashbang */ .chroma .ch { color: #57606a }
/* CommentMultiline */ .chroma .cm { color: #57606a }
/* CommentSingle */ .chroma .c1 { color: #57606a }
/* CommentSpecial */ .chroma .cs { color: #57606a }
/* CommentPreproc */ .chroma .cp { color: #57606a }
/* CommentPreprocFile */ .chroma .cpf { color: #57606a }
/* GenericDeleted */ .chroma .gd { color: #82071e; background-color: #ffebe9 }
/* GenericEmph */ .chroma .ge { color: #1f2328 }
/* GenericInserted */ .chroma .gi { color: #116329; background-color: #dafbe1 }
/* GenericOutput */ .chroma .go { color: #1f2328 }
/* GenericUnderline */ .chroma .gl { text-decoration: underline }
/* TextWhitespace */ .chroma .w { color: #ffffff }
`

// ChromaCSSDark is the Chroma syntax highlighting CSS for the "dracula" style.
// Generated via: chroma/v2 formatters/html WriteCSS with styles.Get("dracula").
// Scoped under .phui-theme-dark so it only applies when dark mode is active.
const ChromaCSSDark = `/* Chroma — dracula style (dark mode) */
.phui-theme-dark .chroma .err { }
.phui-theme-dark .chroma .k { color: #ff79c6 }
.phui-theme-dark .chroma .kc { color: #ff79c6 }
.phui-theme-dark .chroma .kd { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .kn { color: #ff79c6 }
.phui-theme-dark .chroma .kp { color: #ff79c6 }
.phui-theme-dark .chroma .kr { color: #ff79c6 }
.phui-theme-dark .chroma .kt { color: #8be9fd }
.phui-theme-dark .chroma .na { color: #50fa7b }
.phui-theme-dark .chroma .nc { color: #50fa7b }
.phui-theme-dark .chroma .nl { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .nt { color: #ff79c6 }
.phui-theme-dark .chroma .nb { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .bp { font-style: italic }
.phui-theme-dark .chroma .nv { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .vc { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .vg { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .vi { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .vm { color: #8be9fd; font-style: italic }
.phui-theme-dark .chroma .nf { color: #50fa7b }
.phui-theme-dark .chroma .fm { color: #50fa7b }
.phui-theme-dark .chroma .s { color: #f1fa8c }
.phui-theme-dark .chroma .sa { color: #f1fa8c }
.phui-theme-dark .chroma .sb { color: #f1fa8c }
.phui-theme-dark .chroma .sc { color: #f1fa8c }
.phui-theme-dark .chroma .dl { color: #f1fa8c }
.phui-theme-dark .chroma .sd { color: #f1fa8c }
.phui-theme-dark .chroma .s2 { color: #f1fa8c }
.phui-theme-dark .chroma .se { color: #f1fa8c }
.phui-theme-dark .chroma .sh { color: #f1fa8c }
.phui-theme-dark .chroma .si { color: #f1fa8c }
.phui-theme-dark .chroma .sx { color: #f1fa8c }
.phui-theme-dark .chroma .sr { color: #f1fa8c }
.phui-theme-dark .chroma .s1 { color: #f1fa8c }
.phui-theme-dark .chroma .ss { color: #f1fa8c }
.phui-theme-dark .chroma .m { color: #bd93f9 }
.phui-theme-dark .chroma .mb { color: #bd93f9 }
.phui-theme-dark .chroma .mf { color: #bd93f9 }
.phui-theme-dark .chroma .mh { color: #bd93f9 }
.phui-theme-dark .chroma .mi { color: #bd93f9 }
.phui-theme-dark .chroma .il { color: #bd93f9 }
.phui-theme-dark .chroma .mo { color: #bd93f9 }
.phui-theme-dark .chroma .o { color: #ff79c6 }
.phui-theme-dark .chroma .ow { color: #ff79c6 }
.phui-theme-dark .chroma .c { color: #6272a4 }
.phui-theme-dark .chroma .ch { color: #6272a4 }
.phui-theme-dark .chroma .cm { color: #6272a4 }
.phui-theme-dark .chroma .c1 { color: #6272a4 }
.phui-theme-dark .chroma .cs { color: #6272a4 }
.phui-theme-dark .chroma .cp { color: #ff79c6 }
.phui-theme-dark .chroma .cpf { color: #ff79c6 }
.phui-theme-dark .chroma .gd { color: #ff5555 }
.phui-theme-dark .chroma .ge { text-decoration: underline }
.phui-theme-dark .chroma .gh { font-weight: bold }
.phui-theme-dark .chroma .gi { color: #50fa7b; font-weight: bold }
.phui-theme-dark .chroma .go { color: #44475a }
.phui-theme-dark .chroma .gu { font-weight: bold }
.phui-theme-dark .chroma .gl { text-decoration: underline }
.phui-theme-dark .chroma .no { color: #bd93f9 }
.phui-theme-dark .chroma .nd { color: #bd93f9 }
.phui-theme-dark .chroma .ni { color: #f8f8f2 }
.phui-theme-dark .chroma .nn { color: #f8f8f2 }
.phui-theme-dark .chroma .nx { color: #f8f8f2 }
.phui-theme-dark .chroma .p { color: #f8f8f2 }
.phui-theme-dark .chroma .w { color: #f8f8f2 }
`
