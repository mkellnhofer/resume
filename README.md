# Resume

A small static web page generator for my resume.

See: https://mkellnhofer.github.io/resume

![Screenshot](screenshot-light.png#gh-light-mode-only "Screenshot")
![Screenshot](screenshot-dark.png#gh-dark-mode-only "Screenshot")

## Features

- responsive (layouts for desktop, tablet and mobile)
- localizable (a separate page for each language)
- supports "light" and "dark" color schemes 

## Getting started

To get started you need Go.

See: https://golang.org/doc/install

## Configuration

YAML files are used to define resume contents. Each localization has a specific set of files.
Currently, there are files for English and German.

Directory structure:

```
.
+-- data                            <- Resume specific data
|   +-- resume-[LANG-CODE].yaml
+-- strings                         <- General strings
|   +-- strings-[LANG-CODE].yaml
+-- web                             <- Output directory
    +-- flag                        <- Localization icons
    |   +-- [LANG-CODE].svg
    +-- ico                         <- Link icons
    |   +-- address.svg
    |   +-- email.svg
    |   +-- [ICO].svg
    +-- profile-picture.png         <- Profile picture
```

(`LANG-CODE` must be a string with [a-z] and length of 2. All files which don't follow the
convention are ignored.)

If you don't need a localization for a language, just delete the resume YAML file.

### Resume Structure

See: [data/resume-en.yaml](data/resume-en.yaml)

Color dots are used in the skills section to refer to work experience or education.

Following colors are available:

- red
- orange
- yellow
- lightgreen
- green
- lightblue
- blue
- grey (used for self education)

The progress of a major skills can only be [0, 10, ..., 90, 100].

## Usage

Just run:

```
go run resume.go
```

(Currently, there are no options or arguments.)

## Hosting

[GitHub Pages](https://pages.github.com) provides a simple way to host static web pages. You might
want to take a look at this. &#x1F642;

If you want to use a separate branch for the web content, the following guide might be helpful.

https://gist.github.com/cobyism/4730490

## Copyright and License

Copyright Matthias Kellnhofer. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in
compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is
distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing permissions and limitations under the
License.