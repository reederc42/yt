# yt

`yt` is a command line tool and package for manipulating YAML configuration 
files that adds inheritance and components to YAML.

The intent is to make complicated configuration clearer and more maintainable.

In tools such as helm, multiple configuration objects are provided in order. The
leftmost is overridden by the rightmost. This pattern is only enforced by
convention, and can make complex configurations less clear to debug.

`yt` solves this problem by supplying such tools with only one
configuration object, that internally refers to other objects.

## 1 Inheritance

Inheritance is used to enrich or replace parts of a super-object by a
sub-object.

This can be used to specialize configuration objects for different but broadly
similar use cases.

### 1.1 Rules

Inheritance is implemented by two concepts: merge and orthogonal merge.

#### 1.1.1 Merge

`Merge` is overriding elements of a source object by the corresponding elements
of a destination object.

#### 1.1.2 Orthogonal Merge

`Orthogonal Merge` is merging distinct parts of two objects into one object.

Orthogonal merging enables simple multiple inheritance.

## 2 Components

Components are used to provide reuse of sections of configuration objects.

This can be used to maintain common elements across many configuration objects.

Components help break problems that occur with multiple inheritance.
