Reflex Project 
contributors:
hadi tighsazan-402107033
maryam mehdizade-402100526


Project Description

In this project, We worked with the Xray-Core codebase and studied its internal architecture. Xray-Core is a modular proxy platform written in Go that supports multiple protocols such as VMess, VLESS, Trojan, Shadowsocks, SOCKS, HTTP, WireGuard, and others.


Code Structure Overview

We analyzed the following main components:
	•	core/ – Contains the main engine of Xray, configuration loading, and runtime management.
	•	proxy/ – Includes protocol implementations such as VMess, VLESS, Trojan, Shadowsocks, and Reflex.
	•	proxy/reflex/ – Contains the Reflex protocol implementation, including:
	•	handshake logic
	•	codec (packet detection and formatting)
	•	tunnel and session management
	•	inbound and outbound handlers
	•	crypto policy and replay protection
	•	transport/ – Manages network transport layers (TCP, WebSocket, gRPC, TLS, Reality, etc.).
	•	common/ – Utility functions such as cryptography, buffers, protocol helpers, and error handling.
	•	app/ – Application-level services like routing, DNS, logging, policy, stats, and dispatcher.
	•	main/ – Entry point of the program and command-line handling.

We focused especially on the Reflex implementation inside:
proxy/reflex/
We studied how:
	•	The handshake is performed
	•	Encryption and replay protection are handled
	•	Sessions and tunnels are established
	•	Inbound and outbound connections are processed

How to Run the Project
0) Go to the project folder + build Xray
1) Pick ONE UUID (must match on server + client)
2) Create the SERVER config (Reflex inbound)
3) Run the SERVER
4) Create the CLIENT config (SOCKS inbound + Reflex outbound)
5) Run the CLIENT
6) Test that it works


Make sure the configuration file is properly defined for the Reflex protocol if you are testing it.

Problems and Solutions

1. Build Errors (GOPATH / Environment Issues)

Problem:
I encountered errors related to GOPATH and Go environment variables.

Solution:
I checked the Go environment, Then I corrected the PATH and GOPATH configuration to ensure Go could properly resolve dependencies.

2. Dependency and Module Issues

Problem:
Some modules were not resolving correctly during the build process.

Solution:
I used:go mod tidy
to clean and download the required dependencies.

3. Understanding Reflex Architecture

Problem:
The Reflex protocol implementation is spread across multiple folders (handshake, codec, tunnel, inbound, outbound). It was initially difficult to understand how all parts connect together.