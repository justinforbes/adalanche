# Analysis

This is dead simple - everything you've collected should be in the data directory, either in the main folder or in subfolders.

Whatever resides there and Adalanche understands is automatically loaded, correlated and used. It's totally magic.

IMPORTANT: If you're doing multi domain or forest analysis, place all AD object and GPO files for each domain in their own subfolder, so Adalanche can figure out what to merge and what NOT to merge. When dumping just point to a new <code>--datapath</code> for collection run (example: <code>adalanche  --datapath=data/domain1 collect activedirectory --domain=domain1</code> or let Adalanche figure it out itself)

These extensions are recognized:
- .localmachine.json - Windows collector data
- .gpodata.json - Active Directory GPO data
- .objects.msgp.lz4 - Active Directory object/schema data in MsgPack format (LZ4 compressed)

Then analyze the data and launch your browser:

<code>adalanche analyze</code>

There are some options here as well - try <code>adalanche analyze --help</code>

# User Interface

<img src="images/welcome.png" width="80%">

When launched, you'll see some  statistics on what's loaded into memory and how many edges are detected between objects. Don't worry, Adalanche can handle millions of objects and edges, if you have enough RAM ;)

The pre-loaded query allows you to see who can pwn "Administrators", "Domain Admins" and "Enterprise Admins". Query target nodes are marked with RED. 

Press the "analyze" button in the query interface to get the results displayed. If you get a lot of objects on this one, congratz, you're running a pwnshop.

Depending on whether you're over or underwhelmed by the results, you can do adjustments or other searches.

#### Pre-defined searches

To ease the learning experience, there are a number of pre-defined queries built into Adalanche. You access these by pressing the "AQL queries" button, and choosing one. This should give you some idea of how to do queries, but see the dedicated page on that.

#### Analysis Options

<img src="images/analysis-options.png" width="50%">

If your query returns more than 2500 objects (default), Adalanche will limit the output and give you the results that approximately fit within the limit. This limitation is because it has the potential to crash your browser, and is not an Adalanche restriction - feel free to adjust as needed.

The "Prune island nodes" option removes nodes that have no edges from the results.

Analysis depth allows you do limit how many edges from the target selection is returned. Setting this to 0 will only result in the query targets (don't prune islands here, otherwise you'll get nothing), setting it to 1 results on only neighbouring edges to be returned. Quite useful if you get too much data back, blank is no restrictions.

Max outgoing edges limits how many outgoing edges are allowed from an object, and can help limit results for groups and objects that have many assignments. Adalanche tries it best to limit this in a logical way.

Each edge has a probability of success, and you can limit the graph by choosing the minimum probability per edge or overall/accumulated probability of current edge.

#### Edges

Press the "Edges" tab to allow you to do edge based filtering. 

<img src="images/analysis-methods.png" width="50%">

FML is not the usual abbreviation, but represents First, Middle and Last. Disabling the "Middle" selector, will also prevent "Last" in the results, unless it's picked up as the "First" due to the way the search is done.

#### Nodes

<img src="images/object-types.png" width="50%">

This works the same way as the "Edges" limiter above.

#### LDAP query pop-out
When you press the "LDAP Query" tab on the bottom portion of the page, and you get the search interface:

<img src="images/ldap-query.png" width="50%">

You enter a query for things you want to search for, with the "start query" setting your targets. Optionally you can also add a secondary query the following nodes must match. If you put a filter in the "end query" then nodes not matching this will be removed from the outer objects (end of graph).

### Operational theory

Adalanche works a bit differently than other tools, as it dumps everything it can from an Active Directory server, which it then saves to a highly compressed binary cache files for later use. This dump can be done by any unprivileged user, unless the Active Directory has been hardened to prevent this (rare).

If you collect GPOs I recommend using a Domain Admin account, as GPOs are often restricted to apply only to certain computers, and regular users can't read the files. This will limit the results that could have been gathered from GPOs.

The analysis phase is done on all collected data files, so you do not have to be connected to the systems when doing analysis. This way you can explore different scenarios, and ask questions not easily answered otherwise.

### Analysis / Visualization
The tool works like an interactive map in your browser, and defaults to a ldap search query that shows you how to become "Domain Admin" or "Enterprise Admin" (i.e. member of said group or takeover of an account which is either a direct or indirect member of these groups).

### LDAP queries
The tool has its own LDAP query parser, and makes it easy to search for other objects to take over, by using a familiar search language.

**The queries support:**
- case insensitive matching for all attribute names
- checking whether an attribute exists using asterisk syntax (member=*)
- case insensitive matching for string values using equality (=)
- integer comparison using <, <=, > and >= operators
- glob search using equality if search value includes ? or *
- case sensitive regexp search using equality if search value is enclosed in forward slashes: (name=/^Sir.*Mix.*lot$/ (can be made case insensitive with /(?i)pattern/ flags, see https://github.com/google/re2/wiki/Syntax)
- extensible match: 1.2.840.113556.1.4.803 (you can also use :and:) [LDAP_MATCHING_RULE_BIT_AND](https://ldapwiki.com/wiki/LDAP_MATCHING_RULE_BIT_AND) 
- extensible match: 1.2.840.113556.1.4.804 (you can also use :or:) [LDAP_MATCHING_RULE_BIT_OR](https://ldapwiki.com/wiki/LDAP_MATCHING_RULE_BIT_OR) 
- extensible match: 1.2.840.113556.1.4.1941 (you can also use :dnchain:) [LDAP_MATCHING_RULE_IN_CHAIN](https://ldapwiki.com/wiki/LDAP_MATCHING_RULE_IN_CHAIN) 
- custom extensible match: count - returns number of attribute values (member:count:>20 gives groups with more members than 20)
- custom extensible match: length - matches on length of attribute values (name:length:>20 gives you objects with long names)
- custom extensible match: since - parses the attribute as a timestamp and your value as a duration - pwdLastSet:since:<-6Y5M4D3h2m1s (pawLastSet is less than the time 6 years, 5 months, 4 days, 3 hours, 2 minutes and 1 second ago - or just pass an integer that represents seconds directly)
- synthetic attribute: _limit (_limit=10) returns true on the first 10 hits, false on the rest giving you a max output of 10 items
- synthetic attribute: _random100 (_random100<10) allows you to return a random percentage of results (&(type=Person)(_random100<1)) gives you 1% of users
- synthetic attribute: out - allows you to select objects based on what they can pwn *directly* (&(type=Group)(_canpwn=ResetPassword)) gives you all groups that are assigned the reset password right
- synthetic attribute: in - allows you to select objects based on how they can be pwned *directly* (&(type=Person)(_pwnable=ResetPassword)) gives you all users that can have their password reset
- glob matching on the attribute name - searching for (*name=something) is possible - also just * to search all attributes
- custom extensible match: timediff - allows you to search for accounts not in use or password changes relative to other attributes - e.g. lastLogonTimestamp:timediff(pwdLastSet):>6M finds all objects where the lastLogonTimestamp is 6 months or more recent than pwdLastSet
- custom extensible match: caseExactMatch - switches text searches (exact, glob) to case sensitive mode

## Detectors and what they mean

This list is not exhaustive.

| Detector | Explanation |
| -------- | ----------- |
| ACLContainsDeny | This flag simply indicates that the ACL contains a deny entry, possibly making other detections false positives. You can check effective permissions directly on the AD with the Security tab |
| AddMember | The entity can change members to the group via the Member attribute |
| AddMemberGroupAttr | The entity can change members to the group via the Member attribute (the set also contains the Is-Member-of-DL attribute, but you can't write to that) |
| AddSelfMember| The entity can add or remove itself to the list of members |
| AdminSDHolderOverwriteACL | The entity will get it's ACL overwritten by the one on the AdminADHolder object periodically |
| AllExtendedRights | The entity has all extended rights on the object |
| CertificateEnroll | The entity is allowed to enroll into this certificate template. That does not mean it's published on a CA server where you're alloed to do enrollment though |
| ComputerAffectedByGPO | The computer object is potentially affected by this GPO. If filtering is in use there will be false positives |
| CreateAnyObject | Permission in ACL allows entity to create any kind of objects in the container |
| CreateComputer | Permission in ACL allows entity to create computer objects in the container |
| CreateGroup | Permission in ACL allows entity to create group objects in the container |
| CreateUser | Permission in ACL allows entity to create user objects in the container |
| DCReplicationGetChanges | You can sync non-confidential data from the DCs |
| DCReplicationSyncronize | You can trigger a sync between DCs |
| DCsync | If both Changes and ChangesAll is set, you can DCsync - so this flag is an AND or the two others |
| DeleteChildrenTarget | Permission in ACL allows entity to delete all children via the DELETE_CHILD permission on the parent |
| DeleteObject | Permission in ACL allows entity to delete any kind objects in the container |
| DSReplicationGetChangesAll | You can sync confidential data from the DCs (hashes!). Requires DCReplicationGetChanges! |
| GenericAll | The entity has GenericAll permissions on the object, which means more or less the same as "Owns" |
| GPOMachineConfigPartOfGPO | Experimental |
| GPOUserConfigPartOfGPO | Experimental |
| HasAutoAdminLogonCredentials | The object is set to auto login using the entitys credentials which is stored in plain text in the registry for any user to read |
| HasMSA | |
| HasServiceAccountCredentials | The object uses the entitys credentials for a locally installed service, and can be extracted if you pwn the machine |
| HasSPN | The entity has a SPN, and can be kerberoasted by any authenticated user |
| HasSPNNoPreauth | The entity has a SPN, and can be kerberoasted by an unauthenticated user |
| LocalAdminRights | The entity has local administrative rights on the object. This is detected via GPOs or the collector module |
| LocalDCOMRights | The entity has the right to use DCOM against the object. This is detected via GPOs or the collector module |
| LocalRDPRights | The entity has the right to RDP to the object. This is detected via GPOs or the collector module. It doesn't mean you pwn the machine, but you can get a session and try to do PrivEsc |
| LocalSessionLastDay | The entity was seen having a session at least once within the last day |
| LocalSessionLastMonth | The entity was seen having a session at least once within the last month |
| LocalSessionLastWeek | The entity was seen having a session at least once within the last week |
| LocalSMSAdmins | The entity has the right to use SCCM Configuration Manager against the object. This is detected via the collector module. It does not mean that everyone are SCCM admins, but some are |
| MachineScript | Same as above, just as either a startup or shutdown script. Detected via GPOs |
| MemberOfGroup | The entity is a member of this group |
| Owns | The entity owns the object, and can do anything it wishes to it |
| ReadLAPSPassword | The entity is allowed to read the plaintext LAPS password in the mS-MCS-AdmPwd attribute |
| ReadMSAPassword | The entity is allowed to read the plaintext password in the object |
| ResetPassword | The ACL allows entity to forcibly reset the user account password without knowing the current password. This is noisy, and will alert at least the user, who then no longer can log in. |
| ScheduledTaskOnUNCPath | The object contains a scheduled task that sits on a UNC path. If you can control the UNC path you can control what gets executed |
| SIDHistoryEquality | The objects SID-History attribute points to this entity, making them equal from a permission point of view |
| TakeOwnership | The entity can make itself the owner |
| WriteAll | The entity is allowed all write operations |
| WriteAllowedToAct | The entity is allowed to write to the ms-DS-Allowed-To-Act-On-Behalf-Of-Other-Identity attribute of the object, so we can get it to accept impersonations what would otherwise not work |
| WriteAltSecurityIdentities | The entity is allowed to write to the Alt-Security-Identities attribute, so you can put your own certificate there and then authenticate as that user (via PKinit or similar) with this certificate |
| WriteAttributeSecurityGUID | The entity can write to the AttributeSecurityGUID. I'm not sure if this will work, but it has the potential to allows you to add an important attribute to a less important attribute set |
| WriteDACL | The entity can write to the DACL, effectively giving it all permissions after granting them |
| WriteExtendedAll | The entity is allowed to do all extended write operations |
| WriteKeyCredentialLink | The entity can write to the msDK-KeyCredentialLink attribute |
| WriteProfilePath | The entity can write to the user profile path of the user |
| WritePropertyAll | The entity can write to any property (same as above, ACL is just a bit different) |
| WriteScriptPath | The entity can write to the script path of the user, giving them instant remote execution when the user logs on |
| WriteSPN | The entity can freely write to the Service-Principal-Name attributes using SETSPN.EXE or similar tools. You can then kerberoast the account |
| WriteValidatedSPN | The entity can do validated writes to the Service-Principal-Name attributes using SETSPN.EXE or similar tools. You can then kerberoast the account |

## Plotting a path in the GUI

There is a right click menu on objects, so you can to searches in the displayed graph. First right click a target:

<img src="images/set-as-target.png" width="50%">

Then find a source to trace from:

<img src="images/route-to-target.png" width="50%">

If there's a connection from source to target, you'll get the entire attack path presented like this:

<img src="images/found-route.png" width="50%">

You can also pick any object on the graph, and perform an inbound or outbound search from it.
