---
title: "Breaking Bridgefy, again"
subtitle: "Adopting libsignal is not enough"
leadtext: "We analysed the Bridgefy messaging application and found that your private messages are *not* safe."
---

{{< quanda question="What is Bridgefy?" >}}
Bridgefy is a messaging application that uses Bluetooth to transmit messages, so that no Internet is required.
Its developers ([1](https://twitter.com/bridgefy/status/1356750830955884552), [2](https://twitter.com/bridgefy/status/1371507779299590144), [3](https://twitter.com/bridgefy/status/1356680753338318859), [4](https://twitter.com/bridgefy/status/1359200080700600322), [5](https://twitter.com/bridgefy/status/1197191632665415686), [6](https://twitter.com/bridgefy/status/1216473058753597453), [7](https://twitter.com/bridgefy/status/1268015807252004864)) and others ([Reuters](https://www.reuters.com/article/amp/idUSKBN2A22H0), [Forbes](https://web.archive.org/web/20200411154603/https://www.forbes.com/sites/johnkoetsier/2019/09/02/hong-kong-protestors-using-mesh-messaging-app-china-cant-block-usage-up-3685/)) have advertised it for use in areas witnessing large-scale protests and often violent confrontations between protesters and agents of the state.
After [a security analysis in August 2020](https://martinralbrecht.wordpress.com/2020/08/24/mesh-messaging-in-large-scale-protests-breaking-bridgefy/) by Martin R. Albrecht, Jorge Blasco, Rikke Bjerg Jensen, and Lenka Mareková reported severe vulnerabilities, the Bridgefy developers adopted the Signal protocol.
The Bridgefy developers then continued to advertise their application as being suitable for use by higher-risk users, for example during the Russian invasion of Ukraine in 2022 ([8](https://web.archive.org/web/20220302102610/https://twitter.com/bridgefy/status/1495805167920365574), [9](https://web.archive.org/web/20220302102616/https://twitter.com/bridgefy/status/1496981656867016704), [10](https://web.archive.org/web/20220302102619/https://twitter.com/bridgefy/status/1496876534732398593), [11](https://web.archive.org/web/20220302102345/https://twitter.com/bridgefy/status/1498708149083201544)).
{{< /quanda >}}

{{< quanda question="Results" >}}
In this work, we analyse the revised security architecture of Bridgefy and report several vulnerabilities.
The headline news is that we found a practical attack, with a proof-of-concept implementation, that breaks confidentiality of libsignal-protected private messages and succeeds with a probability of about 50%.
Please note that **this attack in no way threatens Signal or libsignal** but attacks how Bridgefy uses it.

Our other findings include:
1. Bridgefy users can still be tracked.
1. Broadcast messages remain unauthenticated; an attacker can exploit this to mount impersonation attacks.
1. The protocol remains susceptible to an attacker in the middle. While such an attack is now limited to the first exchange between a pair of users (i.e., it abuses a "trust on first use" or TOFU assumption) we note that Bridgefy offers users no option to verify the public keys of their contacts.
1. Any nodes in the network that receive a single carefully crafted message become unable to participate in further network communication.
1. The broadcast encryption mechanism employed by the Bridgefy SDK is susceptible to a ciphertext-only attack with the assumption of plaintexts from a small domain. The Bridgefy messenger not affected by this.
{{< /quanda >}}

{{< quanda question="Disclosure" >}}
We disclosed our first vulnerabilities to Bridgefy in May 2021.
According to the developers, the vulnerability allowing an attacker to read encrypted messages was fixed on 14 August 2021.
The disclosure of our attacks on the broadcast encryption followed in September 2021.
We asked the developers to comment on the remediation progress in early February 2022, however, at the time of writing the state of the remediation remains unclear.

We recommend that users avoid Bridgefy until its developers have committed to regular public security audits by respected third party auditors.
{{< /quanda >}}

{{< quanda question="Paper" >}}
You can find more details in [our research paper](./breaking-bridgefy-again.pdf).
{{< /quanda >}}

{{< quanda question="Demo" >}}
A video demo of the TOCTOU attack can be found [here on Twitter](https://twitter.com/eikendev/status/1427542406262575105).
{{< /quanda >}}

{{< quanda question="Exploit Code" >}}
The source code for our attacks is available [here on GitHub](https://github.com/eikendev/breaking-bridgefy-again).
{{< /quanda >}}

{{< quanda question="Team" id="contact" >}}
We are academic researchers from ETH Zurich and Royal Holloway:
- [Martin R. Albrecht](https://malb.io/) (Information Security Group, Royal Holloway, University of London)
- [Raphael Eikenberg](https://www.eiken.dev/) (Applied Cryptography Group, ETH Zurich)
- [Kenneth G. Paterson](https://inf.ethz.ch/people/person-detail.paterson.html) (Applied Cryptography Group, ETH Zurich)
{{< /quanda >}}
