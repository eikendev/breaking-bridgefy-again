---
title: "Breaking Bridgefy, again"
subtitle: "Adopting libsignal is not enough"
leadtext: "We analysed the Bridgefy messaging application and found that your private messages are *not* safe."
---

{{< quanda question="What is Bridgefy?" >}}
Bridgefy is a messaging application that uses Bluetooth to transmit messages, so that no Internet is required.
Its developers ([1](https://twitter.com/bridgefy/status/1356750830955884552), [2](https://twitter.com/bridgefy/status/1371507779299590144), [3](https://twitter.com/bridgefy/status/1356680753338318859), [4](https://twitter.com/bridgefy/status/1359200080700600322), [5](https://twitter.com/bridgefy/status/1197191632665415686), [6](https://twitter.com/bridgefy/status/1216473058753597453), [7](https://twitter.com/bridgefy/status/1268015807252004864)) and others ([Reuters](https://www.reuters.com/article/amp/idUSKBN2A22H0), [Forbes](https://web.archive.org/web/20200411154603/https://www.forbes.com/sites/johnkoetsier/2019/09/02/hong-kong-protestors-using-mesh-messaging-app-china-cant-block-usage-up-3685/)) have advertised it for use in areas witnessing large-scale protests and often violent confrontations between protesters and agents of the state.
After [a security analysis in August 2020](https://martinralbrecht.wordpress.com/2020/08/24/mesh-messaging-in-large-scale-protests-breaking-bridgefy/) by Martin R. Albrecht, Jorge Blasco, Rikke Bjerg Jensen, and Lenka Mareková reported severe vulnerabilities, the Bridgefy developers adopted the Signal protocol.
The Bridgefy developers then continued to advertise their application as being suitable for use by higher-risk users.
{{< /quanda >}}

{{< quanda question="Results" >}}
In this work, we analyse the revised security architecture of Bridgefy and report severe vulnerabilities:
1. Bridgefy users can still be tracked.
1. Broadcast messages remain unauthenticated; an attacker can exploit this to mount impersonation attacks.
1. The protocol remains susceptible to an attacker in the middle. While such an attack is now limited to the first exchange between a pair of users (i.e., it abuses a "trust on first use" or TOFU assumption) we note that Bridgefy offers users no option to verify the public keys of their contacts.
1. Any nodes in the network that receive a single carefully crafted message become unable to participate in further network communication.

The headline news is, however, that we have a practical attack, with a proof of concept implementation, that breaks confidentiality of libsignal-protected private messages which succeeds with a probability of about 50%.

**Our attack in no way threatens Signal or libsignal but attacks how Bridgefy uses it.**
{{< /quanda >}}

{{< quanda question="Disclosure" >}}
We disclosed these vulnerabilities to Bridgefy on 21 May 2021 and the vulnerability allowing an attacker to read encrypted messages was fixed on 14 August 2021. However, we recommend that users avoid Bridgefy until its developers have committed to regular public security audits by respected third party auditors.
{{< /quanda >}}

{{< quanda question="Paper" >}}
You can find more details in [our research paper](./breaking-bridgefy-again.pdf).
{{< /quanda >}}

{{< quanda question="Demo" >}}
A video demo of the TOCTOU attack can be found [here on Twitter](https://twitter.com/eikendev/status/1427542406262575105).
{{< /quanda >}}

{{< quanda question="Exploit Code" >}}
The source code for our attacks will be published at a later stage [here on GitHub](https://github.com/eikendev/breaking-bridgefy-again).
{{< /quanda >}}

{{< quanda question="Team" id="contact" >}}
We are academic researchers from ETH Zurich and Royal Holloway:
- [Raphael Eikenberg](https://www.eiken.dev/) (Applied Cryptography Group, ETH Zurich)
- [Martin R. Albrecht](https://malb.io/) (Information Security Group, Royal Holloway, University of London)
- [Kenneth G. Paterson](https://inf.ethz.ch/people/person-detail.paterson.html) (Applied Cryptography Group, ETH Zurich)
{{< /quanda >}}
