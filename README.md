# Markdown Table Roller
## Usage
gotableroller [tablename]

[tablename] is the name of the markdown file that contains the table you want to roll. The table should be in the directory or subdirectory of gotableroller. It may or may not include the filepath.  

### Example
Given the following directory:
  * Items
    * Weapons.md
    * Materials.md
    * CarriedItems.md

You can roll on the Weapons table with the following commands:
```
  gotableroller Weapons 
  gotableroller Weapons.md
  gotableroller Items/Weapons 
  gotableroller weapons 
```

The table file should contain the contents of the table in a markdown list (ordered or unordered)
  
example:
```
  * Sword
  * Axe
```
or:
```
  1. Sword
  2. Axe
```
The tables can contain other tables as links to other markdown file in the same directory or subdirectory. The links should be in the format of `[name](path/to/file.md)`. 

example:
  * CarriedItems.md
```
    * In the backpack is a [weapon](Items/Weapons.md) made of [material](Items/Materials.md)
    * In the pocket is a trinket made of [material](Items/Materials.md)
```