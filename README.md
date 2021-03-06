# yt

`yt` is a command line tool and package for manipulating YAML configuration 
files that adds inheritance and components to YAML.

The intent is to make complicated configuration clearer and more maintainable.

In tools such as helm, multiple configuration documents are provided in order.
The leftmost is overridden by the rightmost. This pattern is only enforced by
convention, and can make complex configurations less clear to debug.

`yt` solves this problem by supplying such tools with only one configuration
document, that internally refers to other documents.

## 1 Inheritance

Inheritance is used to enrich or replace parts of a super-document by a
sub-document.

This can be used to specialize configuration documents for different but broadly
similar use cases.

### 1.1 Rules

Inheritance is implemented by two concepts: merge and orthogonal merge.

A single document may not inherit directly from multiple documents: each parent
document, even from the same file, must be listed separately, and are subject to
the rules of orthogonal merging.

#### 1.1.1 Merge

`Merge` is overriding elements of a source document by the corresponding
elements of a destination document.

#### 1.1.2 Orthogonal Merge

`Orthogonal Merge` is merging distinct parts of two documents into one document.

Orthogonal merging enables simple multiple inheritance.

##### 1.1.2.1 Logic

Logical description of orthogonal merge.

Assumes two maps, `left` and `right`

Inputs:

a) `key in left`

b) `left[key] is map`

c) `right[key] is map`

Outputs:

Don't panic: `1`

```bash
ab  | 00 01 10 11
----|------------
c 0 |  1  1  0  0
  1 |  1  1  0  1
```

Pseudocode:

```bash
if (key in left && left[key] is map && right[key] is map) || key not in left:
  temp := (key in left) left[key] : nil
  recurse
else:
  panic
```

### 1.2 Syntax

A document may inherit from any UDR. If there is only one document in the parent
file, and the entire file is inherited, the quotes may be elided.

The input is processed line-by-line. Each parent is found by matching this regex 
to each line: `^#&inherits (.*)`

Adding another `#` will disable the inheritance: `##&inherits`

### 1.2.1 Example

`#&inherits filename.yaml`

`#&inherits 'filename.yaml'`

`#&inherits 'filename.yaml'.`

The above lines are equivalent.

## 2 Components

Components are used to provide reuse of sections of configuration documents.

This can be used to maintain common elements across many configuration
documents.

Components help break problems that occur with multiple inheritance.

## 3 Universal Document Reference

UDR may refer to any document from any source.

The first element may be the source document. It can be a http URI. It must be
wrapped in single `'` or double `"` quotes, unless the file contains a single
document and the entire document is used; i.e., if the UDR does not start with 
`'`, `"`, or `.`, the entire string is considered to be the name of a file. 

After the file reference, syntax matches a subset of the `jq` query language.

Specifically, UDR is based on this [path syntax](https://github.com/tidwall/gjson#path-syntax),
prefixed by a file reference.

### 3.1 Multiple Documents

When multiple documents are defined in a file, a single document is referred to
by its position, indexed at 0.

When multiple documents are defined in a file, any UDR referring to that file
must include a specific document.

### 3.2 Examples

| Snippet | Description |
|---|---|
| `#& 'filename.yaml'.` | Sets key as document(s) defined in `filename.yaml` |
| `#& .` | Disallowed: attempts to set the key as the entire current document, a loop |
| `#& '.'.foo` | Refers to the document named 'foo' in the current file |
| `#& .foo` | Refers to the element 'foo' in the current document |

### 3.3 Octothorpe-Ampersand

The octothorpe-ampersand (or pound-sand) is a valid comment in YAML 1.2. With 
`yt`, its use is extended to mean the definition of a parent or query.

## 4 Templates

Beyond the base go template functions, `yaml` and `indent` are defined; with these functions, until otherwise noted, components are implemented with templates.

## 5 Building

### 5.1 Testing

From project root, run `go test ./...`

### 5.2 Installation

`go get github.com/reederc42/yt/cmd/yt`

## 6 Queries & Inserts

The `query` option identifies an element of the provided document.

If `insert` is defined, the identified element (if it exists) is replaed by the value of insert.
