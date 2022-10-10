# Handlebars plugin

This plugin is a wrapper for the [Handlebars for golang](https://github.com/aymerick/raymond) library.

## Installation

Follow the [instructions](https://docs.halon.io/manual/comp_install.html#installation) in our manual to add our package repository and then run the below command.

### Ubuntu

```
apt-get install halon-extras-handlebars
```

### RHEL

```
yum install halon-extras-handlebars
```

## Exported functions

These functions needs to be [imported](https://docs.halon.io/hsl/structures.html#import) from the `extras://handlebars` module path.

### handlebars(template, context)

Render a handlebars template.

**Params**

- template `string` - The template
- context `array` - The context

**Returns**

A successful render will return an associative array with a `result` key that contains the HTML. On error an associative array with a `error` key will be provided.

**Example**

```
import { handlebars } from "extras://handlebars";

$template = ''<div>
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>
'';

echo handlebars($template, [
    "title" => "This is my title",
    "body" =>  "This is my body"
]);
```