# Input Cheatsheet

This cheatsheet is a quick reference for the input data injected into the templates by GenZ.

## Struct

The struct is injected into the template as `Struct`.
It contains the following fields:
- Type (`Type`): The type of the struct.
  - Name (`string`): The name of the type from outside the package. Example `uuid.UUID` or `time.Time`.
  - InternalName (`string`): The name of the type from inside the package. Example `UUID` or `Time`.
- Attributes (`[]Attribute`): List of attributes inside the struct
  - Name (`string`): The name of the attribute.
  - Type (`Type`): The type of the attribute.
    - Name (`string`): The name of the type from outside the package. Example `uuid.UUID` or `time.Time`.
    - InternalName (`string`): The name of the type from inside the package. Example `UUID` or `Time`.
  - Comments (`[]string`): List of comments above the attribute
- Methods (`[]Method`): list of methods associated with the struct.
  - Name (`string`): The name of the method.
  - Params (`[]Type`): List of function paramaters.
    - Name (`string`): The name of the type from outside the package. Example `uuid.UUID` or `time.Time`.
    - InternalName (`string`): The name of the type from inside the package. Example `UUID` or `Time`.
  - Returns (`[]Type`): List of function return values.
    - Name (`string`): The name of the type from outside the package. Example `uuid.UUID` or `time.Time`.
    - InternalName (`string`): The name of the type from inside the package. Example `UUID` or `Time`.
  - Comments (`[]string`): List of comments above the method
  - IsPointerReceiver (`bool`): Whether the method is a pointer receiver or not.
  - IsExported (`bool`): Whether the method is exported or not.
  - Comments (`[]string`): List of comments above the method