# TODO
- [X] Define how to build/use the **set index**: It should be based on the set size defined by the user. Examples of set sizes include: 2, 4, 8, 16, etc. To use module strategy based on the key provided for the user could be an option, but it requires that the key would be an int, that is not possible because there is a business restriction _The keys will always be primitives_. It means it could be an int, string, bool, or something else.
    - [X] Fix the error in the hashKeyToInt function, which returns the same result when inputs of different data types have the same value.

- [X] Implement PUT operation - LRU.
- [ ] Implement PUT operation (MRU algo pending).
- [ ] Implement GET operation.
- [ ] Implement DELETE operation.
- [ ] Implement LISTALL operation.
- [ ] Implement LRU algo.
- [ ] Implement MRU algo.
- [ ] Create README file.

# To keep in mind
- The client interface should be type-safe for keys and values and allow for both the keys and value types to be specified at initialisation time. The keys will always be primitives, and the values could be primitives or other custom types.
- Include some documentation which describes conceptually how your cache works internally and what kinds of use cases it may be suitable for, what it's strengths and weaknesses are, etc.
Document your design in a markdown or pdf file. **Use diagrams** where appropriate. The maximum length of your submission should be 3 A4 pages.

