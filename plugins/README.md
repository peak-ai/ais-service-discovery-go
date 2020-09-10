# Plugins

Service Discovery is built on top of a plug-in architecture. You can simply satisfy the interfaces provided to write your own implementations and integrations with other cloud proviers, other libraries, etc.

This package contains two sub-packages:

1. Core


## Core
Core contains plugins that are used in the wild and are considered stable.


## Incubator
Plugins which are actively being developed and tested against. Which can be used, but with some caution.


## Contributing a plugin
Submit a PR with a plugin, try to make it something useful more broadly, and not overly specific to your own situation. For example, a Kafka adapter would be really useful, or a Google Cloud back-end. However, a locator adapter for a really out of date version of zookeeper that you just so happen to use, might receive pushback, unless there's a case that other users maybe feel this to be useful.

If you don't want to go ahead and do the implementation before floating the idea, go ahead and raise a ticket as a feature request, and we'll see if others feel the same way. If someone has already, upvote it and this will help us to triage.

