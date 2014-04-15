[![Build Status](https://travis-ci.org/dotdoom/goxmpp.png?branch=master)](https://travis-ci.org/dotdoom/goxmpp)

NOTE: XMPP Session Establishment, described in outdated RFC 3921, is no longer a part of standard: [RFC 6121: Functional Summary](http://xmpp.org/rfcs/rfc6121.html#intro-summary).
However since this is a widely used feature, it is implemented in goxmpp.

## [XMPP Standard Compliance](http://xmpp.org/xmpp-protocols/rfcs/)

```
RFC 6120: XMPP CORE                       In progress
  XML Streams                             Yes
  STARTTLS Negotiation                    Yes
  SASL Negotiation                        No
  Resource Binding                        Yes
  XML Stanzas                             In progress
RFC 6121: XMPP IM                         In progress
  Roster Management:                      In progress
RFC 3921 (superseded by RFC 6121):        Rejected
  Session Establishment:                  Yes
```

## [XMPP Extensions (XEPs)](http://xmpp.org/xmpp-protocols/xmpp-extensions/)

```
XEP-0012: Last Activity                   No
XEP-0030: Service Discovery               No
xep-0049: Private XML Storage             In progress
XEP-0077: In-Band Registration            No
XEP-0115: Entity Capabilities             In progress
XEP-0138: Stream Compression              Yes
XEP-0199: XMPP Ping                       No
XEP-0202: Entity Time                     No
XEP-0203: Delayed Delivery                No
XEP-0229: Stream Compression with LZW     In progress
```
