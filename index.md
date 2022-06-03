## Welcome to ECM

ECM stands for Encrypted Content Management and represents a set of tools (CLI, UI, web based) which allows to encrypt your content and manage it safely.

### The Challenge
In modern digital world the content matter the most. We should be confident that it is protected, and stored securely. It should be accessible everywhere, at home, on the web, on a road, etc.

> Companies spend millions of dollars on firewalls and secure access devices, and it is money wasted because none of these measures address the weakest link in the security chain: the people who use, administer and operate computer systems

> Kevin Mitnick

We believe that our customers should be in full control of tools, technology, and moreover have clear understanding how the entire process works. Therefore, we outline every single step of ECM workflows and provide full details of its architecture and open-source codebase

### Machinery
> Technology trust is a good thing, but control is a better one.

> Stephane Nappo

ECM provides tools to encrypt your data in your OS, and safely share it among devices or clients through ECM secure server which can be installed on your premises, within your organization, at cloud provider, etc. All of your data is encrypted before put on a wire. All server communications and APIs use HTTPs protocol to exchange your encrypted data. For complete guide of ECM architecture please follow this link.

ECM provides RESTful APIs to exchange your data within your favorie application. Please refer to our APIs documentation for more details.

### Technology
We provide secure solutions to all of your digital content, including login records, passwords, notes, files, and media. We rely on industry standards to encrypt our data. Client can choose among the following available ciphers:
- AES cipher: The Advanced Encryption Standard (AES), also known by its original name Rijndael is a specification for the encryption of electronic data established by the U.S. National Institute of Standards and Technology (NIST) in 2001.
- NaCl cipher: NaCl (pronounced "salt") is a new easy-to-use high-speed software library for network communication, encryption, decryption, signatures, etc.

Our ECM server relies on LetsEncrypt provider for our server certificates, and we rely on Web Assembly (WASM) technology to enforce the same-origin and permissions security policies of the browser.
