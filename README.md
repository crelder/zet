# zet

The following project proposes a personal knowledge management system that is robust, simple, and built to last a
lifetime.

Inspiration is Vim (about 50 years of development by basically one person), Git (Immutable Architecture, simple interface), Zettelkasten by Niklas Luhmann

## Motivation

We live in a knowledge-based society (German: "Wissensgesellschaft"). There are projects like Wikipedia, where
*collective* knowledge about the world is processed and stored. However, what about *personal* knowledge? Is there a way
to equally enrich how we process and store personal knowledge like the project Wikipedia does for collective knowledge?

This project might be the answer.

Read the following instructions to understand how and why this personal knowledge management system works. If you are
convinced, you can install this CLI tool via `go install github.com/crelder/zet/cmd/cli@latest` (prerequisite: you
have [Go](https://go.dev/doc/install) installed)

## Problem statement

This project tries to solve the question, "What to do with a good thought?". So that I, as the user, have a long-term
memory and can

1. find a thought that I know to exist in my personal knowledge management system and
2. discover relevant thoughts which I do not know to exist in my personal knowledge management system and
3. put a thought in the context of previous thoughts, therefore supporting structured thinking

## Proposed solution

A folder in your computer ("/zettel") holds all your thoughts in the form of a) a textfile or b) a scan of a DIN A6 "zettel" (in German "zettel" means "slip of paper"; there the name "zettelkasten" comes from - "kasten" in German means"
box").

Every filename contains the metadata of your thought and has the structure   
`ID - Keywords - References and/or Contexts - Link`.  
E.g. `170212d - Design, Philosophy - werner2011 243, Movie Matrix - 161103f.txt`.  

Everything but the `ID` is optional in the filename.

### Finding thoughts

Via the search by filename in your file browser and the zettel filename, you can rediscover thoughts. Examples are:

**What was the thought again I had during my vacation in June 2017?**

Since the ID contains the date, your thoughts will appear chronologically in the `zettel` folder, and you can search
for "1706" and it will give you your thoughts during June 2017.

![My thoughts in June 2017](https://github.com/crelder/zettelkasten/blob/master/pictures/search-date-june.PNG)

**What are my thoughts regarding the topic 'Entropy'?**

All keywords are in the filename. So searching for this keyword 'Entropy' would give you this list:

![My thoughts in July 2016](https://github.com/crelder/zettelkasten/blob/master/pictures/search-topic-entropie.PNG)

**What were my thoughts from reading the book from Welter 2011?**

The filename can contain bibkeys for referencing literature sources and a page number. These literature sources are
stored in the `references.bib` file of your zettelkasten. The bibkey's format is AUTHORYEAR (in lower case letters) and
optionally a lower case letter. E.g. `welter2011` or `shannon1948c` in case you have several
references of that author in one year.

![Summary of the book the author Welter wrote in 2011](https://github.com/crelder/zettelkasten/blob/master/pictures/search-source-welter.PNG)

**What was the thing about music and scales in the movie Dunkirk?**

The filename optionally contains context, e.g. the name of a person you had a conversation with when you had this
thought or a place where you had this thought or some other form of context like thinking about a movie.

![My thoughts about the Movie Dunkirk](https://github.com/crelder/zettelkasten/blob/master/pictures/search-source-dunkirk.PNG)

These examples show how you solve problems 1 and 2 in the problem statement above.

### Support for structured thinking

Link each zettel to a previous zettel by providing an ID at the end of a filename. Since the zettel are arranged in your
folder `zettel` in an associative way (ordered by date, since the ID uses the creation date), you can put them in a
logical order by using the entries in the `index.txt` and the command `zet views`. This will create a
folder `VIEWS/index/` which holds for each thematic topic a list of links to your zettel.

You can now use your file browser to display the logical chain of thoughts.

![Chain of thoughts regarding programming](https://github.com/crelder/zettelkasten/blob/b18913a74bccb2dd8abd035e94b9f69371c21d38/pictures/search-structured-thinking.png)

Here you can see that the zettel with ID "220116s" has the ID "220115p" at the end of the filename. If there is more
than one zettel pointing to a previous zettel, it will put the earliest zettel in the main branch and fork all the
others in the form of a folder, e.g. 01_220202g and 02_220519t. In these folders, the chain of thoughts then continues.

An entry in `index.txt` holds a topic and a list of starting points into line of thought in your zettelkasten. Example entries are `Programming: 220115p` or `Entropy: 170213d, 181124s`.

Run `zet init example` to see a simple example zettelkasten - it also serves as a tutorial.

## Design Philosophy

In this section, I want to explain why the above-described solution looks like this and why it is designed in a way that
it fulfills the following requirements.

It is quite some **work** to insert own thoughts into a personal knowledge management system. If I am not convinced that
my thoughts are well protected or are still accessible after a certain time, I will not insert my thoughts in the first
place into the system. This is because I am unwilling to invest time to write thoughts into a system where I cannot
still benefit from my work after a certain time.

From this central observation the following **requirements** for a potential solution derive.

* **Robustness**: The system should not easily break down. Moreover, if it does, you as the user should be able to
  completely recover its state. The system must also work without the need for maintenance (updates, etc.), therefore
  being easy to maintain.

* **Build to last a lifetime**: The system must still work in e.g. 30 years or during your whole lifetime. This is why
  you should not use special software for your personal knowledge management, because it always bears the threat of a "
  look in" of your thoughts into a specific system. Who will guarantee you that you can still use the software in 15
  years? Even if you could export your thoughts, you might not be able to work with your thoughts again. The reason
  could be, that the data is unstructured or messed up. When going bankrupt or discontinuing a product, software
  companies usually are not very dedicated to transfer users' data to another software solution because, in that state,
  they have other sorrows.

* **Simplicity**: To be robust and last a lifetime, the personal knowledge management system must be simple. In general,
  I see that *knowledge management tools are overrated, and the usage of these tools is underrated*. It is as if you
  would give a layman a Steinway grand piano and expected him to play like a great pianist because he has a great tool.
  This is a fallacy.  
  Most work when doing personal knowledge management involves reading, selecting, discussing, thinking, and writing down
  the thoughts.  
  You do not necessarily need `zet` to access your thoughts or to enter your thoughts. `zet` just gives a bit more
  convenience.  
  The assumptions that this system works for you over a lifetime are:

    1. that you will always have a filesystem with a folder, where you can store text and image files.
    2. that you will always have a search functionality for searching through your filenames, e.g. via the file
       explorer (Windows) or Finder (Mac).
    3. that you have a viewer and editor for textfiles and a viewer for images.

* **Immediate use**: You should not have to enter content in the system for e.g. one year until a useful experience is
  promised. The system should immediately generate a useful output and, therefore, an experience of genuine use (after
  having added e.g. five thoughts or having used it for three days). This constant feedback on the use of the system
  keeps you motivated to input data into the system continuously.

* **Good input**: The outcome quality of a system depends on the input quality. A system must be designed in a way so
  that no "bad" input enters the system. No matter how good the system is, if you enter bad input, the system will have
  bad output - independent of the system's quality: f(shit) = shit. What is meant by "bad input" are thoughts that are
  not your own thoughts and are merely a "copy-paste" of what someone else said or wrote. "Good input" are thoughts
  (even from someone else) which I can write down in consistent prose text using my **own words** and optionally enrich
  it with sketches, arrows, etc. See also
  this [explanation](https://strengejacke.wordpress.com/2007/08/04/lesen-lernen/)
  on how to read.

A consciously chosen limit in this approach to personal knowledge management is that you decide to use only one language
for the filenames.

## Further reading

The full potential of this personal management system is unlocked when combining this tool with other tools for e.g.
synchronizing via the cloud, version control. See the [FAQs](docs/FAQ.md) about these points and further reading.

There are more `zet` features. `zet help` will show you a list of all commands. `zet init example` will download a small
tutorial, where especially the command `zet import` is explained.

## Feedback and Contributing

I am happy to receive comments, thoughts and ideas in
the [discussion section](https://github.com/crelder/zet/discussions/1).

I develop this project in my spare time. Therefore, I am happy about any contributions to this project, be it new
features, bug fixes, typo fixes, documentation improvements, suggestions on how to design the code better, etc.

Just open a new issue ticket or a pull request. I promise to respond, but it might take a while because I am busy with
other stuff.

## Acknowledgments

Gabriel, Silke, Sascha, Mathias, Siegfried