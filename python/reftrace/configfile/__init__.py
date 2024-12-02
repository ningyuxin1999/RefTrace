from dataclasses import dataclass
from ..proto import config_file_pb2

@dataclass(frozen=True)
class Value:
    """Represents a directive value with parameters and closure information."""
    _value: config_file_pb2.DirectiveValue

    @property
    def params(self) -> list[str]:
        """The directive parameters."""
        return list(self._value.params)

    @property
    def in_closure(self) -> bool:
        """Whether the directive is in a closure."""
        return self._value.in_closure

@dataclass(frozen=True)
class NamedOption:
    """Represents a named option with its value."""
    _value: config_file_pb2.NamedOption

    @property
    def line_number(self) -> int:
        """The line number where this option appears."""
        return self._value.line_number

    @property
    def name(self) -> str:
        """The option name."""
        return self._value.name

    @property
    def value(self) -> Value:
        """The option value."""
        return Value(self._value.value)

@dataclass(frozen=True)
class Directive:
    """Represents a directive configuration."""
    _value: config_file_pb2.DirectiveConfig

    @property
    def line_number(self) -> int:
        """The line number where this directive appears."""
        return self._value.line_number

    @property
    def name(self) -> str:
        """The directive name."""
        return self._value.name

    @property
    def options(self) -> list[NamedOption]:
        """The directive options."""
        return [NamedOption(opt) for opt in self._value.options]

    @property
    def value(self) -> Value:
        """The directive value."""
        return Value(self._value.value)

@dataclass(frozen=True)
class NamedScope:
    """Represents a named scope with its directives."""
    _value: config_file_pb2.NamedScope

    @property
    def line_number(self) -> int:
        """The line number where this scope begins."""
        return self._value.line_number

    @property
    def name(self) -> str:
        """The scope name."""
        return self._value.name

    @property
    def directives(self) -> list[Directive]:
        """The directives in this scope."""
        return [Directive(d) for d in self._value.directives]

@dataclass(frozen=True)
class ProcessScope:
    """Represents a process scope containing directives and named scopes."""
    _value: config_file_pb2.ProcessScope

    @property
    def line_number(self) -> int:
        """The line number where this scope begins."""
        return self._value.line_number

    @property
    def directives(self) -> list[Directive]:
        """The directives in this scope."""
        return [Directive(d) for d in self._value.directives]

    @property
    def named_scopes(self) -> list[NamedScope]:
        """The named scopes within this scope."""
        return [NamedScope(s) for s in self._value.named_scopes]

__all__ = [
    'ProcessScope',
    'NamedScope',
    'Directive',
    'NamedOption',
    'Value',
]