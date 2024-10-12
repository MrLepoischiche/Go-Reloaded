# Go-Reloaded
A project made during studies : a program manipulating words and correcting grammar and punctuation mistakes in an English text, written in Go language.

## What exactly does it do?
In short, this program takes a text file as input, and outputs another text file, with all the necessary corrections applied to the inputted text.

For instance, an incorrect text such as `Then he said" Why the long face ? " ,so I said " Why the small D ? " . ` will be corrected to `Then he said "Why the long face?", so I said "Why the small D?".`.
---
This program also supports the use of word manipulation tags, inside brackets.

Available tags include:
- `up`: Uppers words (`I am so (up) good at this game .` -> `I am SO good at this game.`)
- `low`: Lowers words (`Why do I yell so OFTEN (low) ?` -> `Why do I yell so often?`)
- `cap`: Capitalizes words (`No one comes to the father (cap) except through me.` -> `No one comes to the Father except through me.`)
- `bin`: Converts a Binary number into its Decimal representation (`110 (bin) DOLLARS SHRIMP SPECIAL` -> `6 DOLLARS SHRIMP SPECIAL`)
- `hex`: Converts a Hexadecimal number into its Decimal representation (`138D5 (hex)` -> `80085`)

Tags also support the use of numbers, separated from the tag itself with a comma: `(up, 3)`

It specifies the amount of words to modify *before* the tag:
`Monarchy names are stupid. I can be named lord harrington bigwood the third (up, 5) and no one will question it.`
`Monarchy names are stupid. I can be named Lord Harrington Bigwood The Third and no one will question it.`