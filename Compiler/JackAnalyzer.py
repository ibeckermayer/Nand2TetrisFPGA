from typing import TextIO
from enum import Enum


class TokenType(Enum):
    KEYWORD = 0
    SYMBOL = 1
    IDENTIFIER = 2
    INT_CONST = 3
    STRING_CONST = 4


TT = TokenType


class KeyWord(Enum):
    CLASS = 0
    METHOD = 1
    FUNCTION = 2
    CONSTRUCTOR = 3
    INT = 4
    BOOLEAN = 5
    CHAR = 6
    VOID = 7
    VAR = 8
    STATIC = 9
    FIELD = 10
    LET = 11
    DO = 12
    IF = 13
    ELSE = 14
    WHILE = 15
    RETURN = 16
    TRUE = 17
    FALSE = 18
    NULL = 19
    THIS = 20


KW = KeyWord


class JackTokenizer:
    '''
    Takes in a .jack input file and steps through each character (as advance() is called), updating it's internal state to represent the
    input as Jack-language tokens, as specified by the Jack grammar
    '''

    def __init__(self, filename: str):
        self._filename = filename
        self._stream: str = open(self._filename, 'r').read()  # Reads the entire file as a string
        self._streamlen = len(self._stream)
        self._i: int = 0  # Index of the character in the stream we are currently analyzing
        self.has_more_tokens = True
        self.token_type: TokenType
        self._symbol: str  # Becomes a get-able property when token_type == TT.SYMBOL
        self._int_val: int  # Get-able when token_type == TT.INT_CONST
        self._string_val: str  # Get-able when token_type == TT.STRING_CONST
        self._key_word: KeyWord  # Get-able when token_type == TT.KEYWORD
        self._identifier: str  # Get-able when token_type == TT.IDENTIFIER

    def advance(self):
        '''
        Each [non-recursive] call to advance() eats the next token in the _stream and and updates self's internal state
        to reflect the token it just ate. Specifically, it updates token_type and sets whichever of _symbol, _int_val, _string_val,
        _key_word, or _identifier corresponds.
        '''
        if self._i >= self._streamlen:
            self.has_more_tokens = False
            return

        # Skip all whitespace characters
        if self._stream[self._i].isspace():
            self._i += 1
            self.advance()
            return

        # Skip line comments
        if self._stream[self._i] == '/' and self._stream[self._i + 1] == '/':
            self._i += 2
            # Advance until either a newline or EOF
            while self._i < self._streamlen and self._stream[self._i] != '\n':
                self._i += 1
            self._i += 1  # Eat the '\n'
            self.advance()
            return

        # Skip block comments
        if self._stream[self._i] == '/' and self._stream[self._i + 1] == '*':
            self._i += 2
            while self._i < self._streamlen and not (self._stream[self._i] == '*' and
                                                     self._stream[self._i + 1] == '/'):
                self._i += 1
            self._i += 2  # Eat the closing '*/'
            self.advance()
            return

        # Now that comments and whitespace have been skipped, we can determine what type of lexical element we're at
        if self._stream[self._i] in '{}()[].,;+-*/&|<>=~':
            # If we're at a single character symbol
            self.token_type = TT.SYMBOL
            self._symbol = self._stream[self._i]  # Set _symbol to the symbol
            self._i += 1  # Eat the symbol
        elif self._stream[self._i].isnumeric():
            # If we're at an integer constant
            self.token_type = TT.INT_CONST
            j = self._i  # Save our current index
            self._i += 1  # Eat the current char
            # Eat the subsequent numeric values to capture the full constant
            while self._stream[self._i].isnumeric():
                self._i += 1
            self._int_val = int(self._stream[j:self._i])
            if self._int_val > 32767:
                raise RuntimeError(
                    f"Integer constants must be a decimal number in the range of 0 .. 32767. Got {self._int_val}"
                )
        elif self._stream[self._i] == '"':
            # If we're at a string constant
            self.token_type = TT.STRING_CONST
            self._i += 1  # Eat the '"'
            j = self._i  # Save our current index
            while self._i < self._streamlen and self._stream[self._i] != '"':
                # Eat the rest of the string constant
                if self._stream[self._i] == '\n' or elf._stream[self._i] == '\r':
                    raise RuntimeError("String constants cannot contain newline characters")
                self._i += 1
            if self._i >= self._streamlen:
                # Check that we didn't hit EOF without the string constant terminating
                raise RuntimeError("Hit EOF before string constant terminated")
            if self._i == j:
                # TODO: Maybe should support empty strings?
                raise RuntimeError("Encountered unsupported empty string constant")
            self._string_val = self._stream[j:self._i]
            self._i += 1  # Eat the closing '"'
        else:
            # We're at either an identifier or a keyword or an invalid character
            if not (self._stream[self._i].isalpha() or self._stream[self._i] == '_'):
                # Check for invalid character
                raise RuntimeError(f"Encountered invalid character: {self._stream[self._i]}")

            j = self._i  # Save current index
            while self._i < self._streamlen and (self._stream[self._i].isalnum() or
                                                 self._stream[self._i] == '_'):
                # Eat the rest of the characters in this identifier or keyword
                self._i += 1

            if self._stream[j:self._i] in [
                    'class', 'constructor', 'function', 'method', 'field', 'static', 'var', 'int',
                    'char', 'boolean', 'void', 'true', 'false', 'null', 'this', 'let', 'do', 'if',
                    'else', 'while', 'return'
            ]:
                # If this is a keyword
                self.token_type = TT.KEYWORD
                self._key_word = KeyWord(self._stream[j:self._i])
            else:
                # Else this is an identifier
                self.token_type = TT.IDENTIFIER
                self._identifier = self._stream[j:self._i]

    @property
    def symbol(self) -> str:
        if self.token_type != TT.SYMBOL:
            raise RuntimeError(
                "Attempted to access symbol but the current token is not of type TT.SYMBOL")
        return self._symbol

    @property
    def int_val(self) -> int:
        if self.token_type != TT.INT_CONST:
            raise RuntimeError(
                "Attempted to acces int_val but the current token is not of type TT.INT_CONST")
        return self._int_val

    @property
    def string_val(self) -> str:
        if self.token_type != TT.STRING_CONST:
            raise RuntimeError(
                "Attempted to acces string_val but the current token is not of type TT.STRING_CONST"
            )
        return self._string_val

    @property
    def key_word(self) -> KeyWord:
        if self.token_type != TT.KEYWORD:
            raise RuntimeError(
                "Attempted to acces key_word but the current token is not of type TT.KEYWORD")
        return self._key_word

    @property
    def identifier(self) -> str:
        if self.token_type != TT.IDENTIFIER:
            raise RuntimeError(
                "Attempted to acces identifier but the current token is not of type TT.IDENTIFIER")
        return self._identifier
