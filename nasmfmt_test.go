package main

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestProcess_HappyPath(t *testing.T) {
	source := `%include "const.h"

        variable      db      "text"    ; the comment

_label:
    mov ax, cs
    mov       ss, ax
    jmp     _new_label
;; will be  not executed.
jmp $   ;
    _new_label:
mov ax cs
;; section_end _new_label.
        mov es, ax
`
	formatted := `%include "const.h"

variable db "text" ; the comment

_label:
    mov ax, cs
    mov ss, ax
    jmp _new_label
    ;; will be  not executed.
    jmp $ ;
_new_label:
    mov ax cs
    ;; section_end _new_label.
mov es, ax
`
	reader := bufio.NewReader(strings.NewReader(source))
	buffer := bytes.Buffer{}
	process(reader, &buffer)
	assert.Equal(t, formatted, buffer.String())
}

func TestProcess_CommentNotTrimed(t *testing.T) {
	source := `mov es, ax ; comment line  with two sapces.
`
	reader := bufio.NewReader(strings.NewReader(source))
	buffer := bytes.Buffer{}
	process(reader, &buffer)
	// The source will not be formatted
	assert.Equal(t, source, buffer.String())
}

func TestProcess_StringNotTrimed(t *testing.T) {
	source := `variable db "the string contains spaces."
`
	reader := bufio.NewReader(strings.NewReader(source))
	buffer := bytes.Buffer{}
	process(reader, &buffer)
	// The source will not be formatted
	assert.Equal(t, source, buffer.String())
}

func TestProcess_AnotherStringCase(t *testing.T) {
	source := `variable db "the string contains \" \" spaces."
`
	reader := bufio.NewReader(strings.NewReader(source))
	buffer := bytes.Buffer{}
	process(reader, &buffer)
	// The source will not be formatted
	assert.Equal(t, source, buffer.String())
}

func TestProcess_ColonInSuffixComent(t *testing.T) {
	source := `mov ss, ax ; the comment:
mov cs, ax
`
	reader := bufio.NewReader(strings.NewReader(source))
	buffer := bytes.Buffer{}
	process(reader, &buffer)
	// The source will not be formatted
	assert.Equal(t, source, buffer.String())
}

func TestProcess_ColonInLineComent(t *testing.T) {
	source := `;; the comment:
mov cs, ax
`
	reader := bufio.NewReader(strings.NewReader(source))
	buffer := bytes.Buffer{}
	process(reader, &buffer)
	// The source will not be formatted
	assert.Equal(t, source, buffer.String())
}
