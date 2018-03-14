# GetFile

Skill-testing demo for retrieving a remote URL in pieces.

see REQUIREMENTS.md for more details on the requirements.

## Methodology

### approach

The requirements could easily have been met with a functional approach.  No objects were really needed. Objects were
used for the following reasons:

- isolate some of the risky parts to a utility class, which could be upgraded in the future
- try to prove some understanding of the language

@NOTE - no interfaces were defined, as it was felt that this would really have been superfluous.  Interfaces are
considered important OOP and golang structures, but it could not be justified in this tool.

### architecture

#### abstraction

The primary functionality abstraction is placed in the GetFile struct, which can be used to retrieve a file in pieces.

This tool has a method ::GetPieces() which meets the requirements for the tool.  It could use more work:

- it should validate file length versus retrieval request
- It should allow an approach to download a full file, defining the number of pieces, or the piece size.

the cli command is provided using a cli implementation that uses the GetFile struct, with some logging.

#### cli

the cli is built using the urfave/cli, which is just a common cli lib used for arg interpretation, and clean handling

#### logging

pretty general output logging is handled in the cli part of the app.  It is a it stupid and verbose, but the output,
which is run through the common logrus/log library, is limited to cli handler, and can be avoided if the tooling is
needed to be used as a library.

The functionality of the tool available as a library handles logging primarily though error returning.

### testing

Testing is implemented using standard go TDD tools.  As structs/tools were developed, they were written using the
following tests:

- define interface (actually just struct & methods, as we didn't use any interfaces)
- write some tests for the public methods
- flesh out the tests and the functionality

Then tests were used to ensure functionality and prevent regressions.

TESTING is a part of the build process.  There is a make target for testing, which should be used as a part of the
contribution process.
