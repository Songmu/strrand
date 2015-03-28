strrand
=======

[![Build Status](https://travis-ci.org/Songmu/strrand.png?branch=master)][travis]
[![Coverage Status](https://coveralls.io/repos/Songmu/strrand/badge.png?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/strrand?status.svg)](godoc)

[travis]: https://travis-ci.org/Songmu/strrand
[coveralls]: https://coveralls.io/r/Songmu/strrand?branch=master
[license]: https://github.com/Songmu/strrand/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/strrand

## Description

generate random strings from a pattern like regexp

It is golang porting of perl's [String::Random](https://metacpan.org/release/String-Random) (ported only `randregex` interface)

## Synopsis

    str, err := strrand.RandomString(`[1-9]{1,3}.?`)
    fmt.Println(str) // 13h

OO interface

    sr := strrand.New()
    str, err := sr.Generate(`[あ-お]{15,18}`)

Factory Method

    g, err := strrand.New().CreateGenerator(`\d{2,3}-\d{3,4}-\d{3,4}`)
    str1 := g.Generate() // 11-2258-333
    str2 := g.Generate() // 093-0033-3349

## Supported Patterns

Please note that the pattern arguments are not real regular expressions. Only a small subset of regular expression syntax is actually supported. So far, the following regular expression elements are supported:

    \w    Alphanumeric + "_".
    \d    Digits.
    \W    Printable characters other than those in \w.
    \D    Printable characters other than those in \d.
    \s    Whitespaces (whitespace and tab character)
    \S    Ascii characters without whitespaces
    .     Printable characters. (ascii only)
    []    Character classes. (Supported multibyte characters)
    {}    Repetition.
    *     Same as {0,}.
    ?     Same as {0,1}.
    +     Same as {1,}.

## Disclaimer

Seeding is naive and not secure. So, don't use this for creating password and so on.

## Author

[Songmu](https://github.com/Songmu)
