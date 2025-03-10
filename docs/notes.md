# General generation process
thoughts on how this works:
1. load a file with functions
2. parse the file into functions
3. generate the binary with all functions
4. find the function that matches the requested function
5. parse arguments from the cmdLine or input into the desired types
6. pass the arguments to the function

# Goals of this tool are to work in the following ways:
1. recursively search the current and parent folders for:
   1. a ".gogo" folder
   2. a "gogofiles" folder
   3. a ".gogobuild" file, which is a gogo configuration file (To be designed): we follow directions in config
   4. a ".git" folder: we stop if nothing is found
   5. the root of the filesystem: we stop if nothing is found
2. if nothing is found that matches the requested function, it should search the global gogo namespace
3. there must be a way to force using the global version, either using environment variable or the override described below (this should most likely be rarely done)

We intentionally *DO NOT* search the $PWD for gogo files. This is to prevent the user from accidentally running a function that they didn't intend to run.

# Arguments/Flags
All positional arguments to the gogo command are passed to the function as arguments, with the exception of the first argument, 
which is the function name. however, we *can* pass flags to gogo.

to pass these arguments to the actual function, we need to use the -- separator.
so like `gogo function -- --help` would pass `--help` to the function

# How is this unique from magefiles / bake?

*All of this is very subjective!*

* just better: uses cobra/viper for the output binary 
  * which provides shell completions from the output binary (maybe magefile does this, bake def does)
  * because we use viper, the resulting tools can have arguments be passed in from the command line, environment variables, or a configuration file, or any combination of the above
* ! better than mage: due to the fact we specify argument constraints in the function using the context, we can parse this and provide better help/autocomplate/information to the user in the CLI.
* ! better than bake: bake does this as well, however it's function arguments get lost in many stringly typed option calls. gogo keeps the function arguments in the signature of the function, which makes it easier to read and update
* just better: the gogobuild file can be toml, yaml, or pkl. This allows for more flexibility in the configuration file
* ! just better: gogo searches a global namespace for functions, which allows gogo to replace system-wide bash scripts/functions
* ?: one restriction is module/local functions cannot inherently call global functions. (I think this is ok)
* just better: gogo will watch a global namespace folder(s) and automatically re-compile as necessary
* just better: the global gogo configuration file can be in toml, yaml, or pkl
* worse: not as well tested/used/vetted
* worse: not as well documented
* worse: targeting darwin/zsh and linux/bash for now
* worse: volatile at the moment
* way worse: no dependency management, so we can't: 
  * instantiate dependencies that exist in the function signature
  * require other functions to run before the target function

# gogo build modes
* by default, when a .gogo or other folder is found, it generates a binary in that folder called "gg_binary[_timestamp]"
* by default, when file(s) in the local folder are found to have the gogo build tag, it generates a binary in the local folder called "gg_binary[_timestamp]"
* by default, the output names of the binaries are in lower snake case
* by default, the output binary for global functions is ggg and placed in the $GOGO_BIN folder
* a local function may declare that it be built alone, without any other functions included. In this case the binary is named after the function, lower_snake_case. it gets built in the same folder as the file
* a global function may declare that it be built alone, without any other functions included. In this case the binary is named after the function, lower_snake_case. it gets built in the $GOGO_BIN folder
* even when a global function is built alone, it is still included in the global ggg binary

# Autocomplete
The auto-completion for the gogo binary is unique in that:
1. If no arguments are given, the gogo binary determines a list of all possible targets and returns that
2. If a target is given, then:
    1. the shell completion calls gogo with `--autocomplete=<target>` and the gogo binary returns the binary and base command to perform auto-completions against
    2. then the shell completion takes that result, and calls the binary+base command to get the results to display to the user

This means that:
1. The Gogo binary has unique auto-complete scripts, which are not the generated/default ones by cobra
2. each built binary needs to have an auto-generated cobra `cmd.ValidArgsFunction` which returns possible values for auto-completion. This is generated from the function ctx.Argument.Options()

# Binary Output
In the normal mode, we'll generate a single binary for all the functions in the entire directory. 
This binary will have a command for each function, and the function will be called when the command is called.

In the secondary mode, we'll output a single binary for every function in the directory. This is useful for outputting functions
directory into the local namespace, allowing the functions to be called like scripts. All global functions have this done automatically.

This means that the configuration of this is understood before reading the files themselves (driven by arguments to the gogo binary)

# Scenarios

## Listing

#### From outside local gogo directory and no local gogo files
If we run `gogo` from outside a gogo directory, we should see the global functions listed
It should output something like:
```
Global Functions:
  * function1: (4 args) description
  * function2: (3 args) description
  * function3: (0 args) description
```

Let's try that again. By "path" I mean the current "path" in the current line being inputted in the shell. So for example:

Let's say I am providing autocompletion for the "gg" binary.
If I type "gg" and hit <tab>
I want autocompletion to send `gogo --autocomplete=.`
This should return a list of values.

Then, let's say they select one:
"gg firstOption" and hit tab,

I want the autocomplete to send `gogo --autocomplete=firstOption.1`
Then the shell should display the options.

Then, let's say they further refine their selection:
"gg firstOption arg1" and they hit tab

I want the autocomplete to send `gogo --autocomplete=firstOption.2`

# Scenarios

## subdir
It's assumed that the gogo tool is run from within subdir. This tests finding the .gogo folder and running the functions from within.
This also shares the go.mod from the parent director

## subdir unique gomod
Same as subdir, but the subdir has it's own go.mod file

## Testing Snapshots Library
https://github.com/bradleyjkemp/cupaloy

## Project Layout
* Inside the applications that get built, I want them to source 

## GoGo Gadget
I want people to see "gogo gadget <thing>" in configuration. 
Either that or I want "gogo" to call "gadgets". This means that gadgets would need to be go functions, and gogo would be the task runner.
I like "gogo" as the name for the go function runner, b/c it's basically "go"ing "go", with is what it literally does.
But if "gogo" is the go function runner, then "gadgets" would be tasks, and tasks would then be calling "gogo", which is opposite of what I want.

If I keep gogo as the go function runner, and "gadget" for the task runner, I can have a mode in "gogo" that allows you to do:
`gogo gadget <thing>` and it'll just call the gadget version of the thing. Just for fun.

## Variable checking that needs to occur
We need to make sure that no variables are named "help" or "h"
Also need to check that no ctx().var(<name>).Name() contains no space