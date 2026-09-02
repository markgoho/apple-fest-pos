# POS touch UI design guidance, and the young-operator question

Research note. Compiled 2026-09-01 for the Apple Fest POS.

## Question

What published, citable design guidance exists for point-of-sale (POS) touch UIs, and does any of it address young operators (roughly ages 14 to 16) running a till?

## Summary

Two findings, one per part.

**Part A.** Citable guidance with hard numbers exists, but almost none of it is about point of sale. The POS industry publishes essentially no design specifications: of six major vendors, only Clover states a number and it is borrowed from WCAG, and only Oracle publishes figures, which are database capacity limits. The usable numbers come from three places instead: the accessibility standard (WCAG 2.2), the two platform guidelines (Apple, Material), and the touchscreen ergonomics literature. These three disagree by a factor of more than two, and the disagreement is systematic: **Apple's 44 pt and Material's 48 dp are phone-thumb figures, while the kiosk literature converges on 20 mm for a person standing at a screen and touching it with a finger.** A festival till is the second case.

**Part B.** Nothing exists. No standards body, platform owner, POS vendor, or HCI publication issues interface design guidance for operators aged 14 to 16. The touchscreen motor-performance literature has a documented hole at exactly this band: it is dense for ages 3 to 10, dense again from 18 upward, and empty between. The clearest proof is the MTAGIC project, the field's most comprehensive children's-touchscreen research programme, which was scoped for ages 5 to 17 and published only on 5 to 10. The one legal fact that does exist points the other way from what one might expect: 29 CFR 570.34(d) names "Cashiering" as expressly permitted work for 14- and 15-year-olds, and a Scout volunteering at a nonprofit booth is probably outside the FLSA entirely. **The correct design conclusion is that the operator's age is not a sizing input. Their inexperience is a workflow input.** Design the targets for a standing kiosk and design the error recovery for a novice.

## How to read the evidence classes

- **Hard standard** - a normative specification or a regulation. WCAG 2.2, ISO, the Code of Federal Regulations.
- **Platform guidance** - a platform owner tells developers what to do. Apple HIG, Material Design 3. Not normative, but stable and citable.
- **Vendor guidance** - a POS company publishes design rules for its own product.
- **Academic finding** - a peer-reviewed measurement. Cite the study, not the conclusion.
- **Absence of evidence** - searched for, not found. Recorded on purpose.

---

## PART A - POS touch UI guidance

### A1. Apple Human Interface Guidelines

**Class: platform guidance.**

Apple publishes the control-size numbers on the Accessibility page of the Human Interface Guidelines, under "Offer sufficiently sized controls". The current page gives a table with two columns, a *default* size and a *minimum* size. This matters, because the number that circulates in secondary write-ups ("44pt is the Apple minimum") is the **default**, not the minimum.

> "Controls that are too small are hard for many people to interact with and select. Strive to meet the recommended minimum control size for each platform to ensure controls and menus are comfortable for all when tapping and clicking."

| Platform | Default control size | Minimum control size |
| --- | --- | --- |
| iOS, iPadOS | 44x44 pt | 28x28 pt |
| macOS | 28x28 pt | 20x20 pt |
| tvOS | 66x66 pt | 56x56 pt |
| visionOS | 60x60 pt | 28x28 pt |
| watchOS | 44x44 pt | 28x28 pt |

Source: Apple, *Human Interface Guidelines - Accessibility*, <https://developer.apple.com/design/human-interface-guidelines/accessibility> (table quoted verbatim; retrieved 2026-09-01 from the page data endpoint <https://developer.apple.com/tutorials/data/design/human-interface-guidelines/accessibility.json>, because the HTML page is a JavaScript shell).

Apple also gives spacing numbers on the same page, which are more useful for a dense POS grid than the size number alone:

> "Consider spacing between controls as important as size. Include enough padding between elements to reduce the chance that someone taps the wrong control. In general, it works well to add about 12 points of padding around elements that include a bezel. For elements without a bezel, about 24 points of padding works well around the element's visible edges."

Same source.

**Unit note.** On iOS and iPadOS 1 pt equals 1 CSS px in a web view at the default zoom. So 44 pt reads as 44 CSS px, and Apple's 12 pt padding reads as 12 CSS px.

**Apple contradicts itself.** The Buttons page states the 44 pt figure as a floor, not a default:

> "a button needs a hit region of at least 44x44 pt - in visionOS, 60x60 pt - to ensure that people can select it easily."

Source: Apple, *Human Interface Guidelines - Buttons*, <https://developer.apple.com/design/human-interface-guidelines/buttons>.

So the Buttons page says "at least 44x44 pt" while the Accessibility page says the minimum is 28x28 pt and 44x44 pt is the default. The safe reading for a POS is the stricter one: treat 44 pt as a floor and ignore the 28 pt figure.

The Layout page carries no numeric tap-target guidance at all, only the qualitative "Make controls easier to use by providing enough space around them and grouping them in logical sections." <https://developer.apple.com/design/human-interface-guidelines/layout>.

**Not found.** Apple publishes no point-of-sale or kiosk section in the Human Interface Guidelines. Searches of the HIG index found no page for POS, kiosk, retail, or "single-purpose device". The closest published guidance is Apple's Guided Access / Autonomous Single App Mode documentation, which is a device-management feature, not design guidance.

### A2. Material Design 3 and Android

**Class: platform guidance.**

Google states 48dp, and states it in more than one place.

The Material Design 3 page *Foundations - Designing - Structure*, <https://m3.material.io/foundations/designing/structure>, states:

> "For most platforms, consider making touch targets at least 48 x 48dp."

> "In most cases, targets separated by 8dp of space or more promote balanced information density and usability."

and cross-references Apple:

> "iOS recommends 44 x 44dp targets."

**Verification caveat.** m3.material.io is a JavaScript application that will not serve its text to an automated fetcher, so the three quotes above could not be re-verified independently at the raw-HTTP level. They are corroborated by Google's own server-rendered documentation below, which states the same 48dp rule. A reader who needs the 8dp spacing figure specifically should open the m3 page in a browser to confirm it.

The corroborating server-rendered sources:

> "For touch interfaces, we recommend that each interactive UI element have a focusable area, or touch target size, of at least 48dpx48dp. Larger is even better."

Source: Android Developers, *Make apps more accessible*, <https://developer.android.com/guide/topics/ui/accessibility/apps>.

The same page adds an exception that a POS on a tablet does not get to use:

> "For precise input (mouses and trackpads), the touch target can be smaller."

Google's Android Accessibility Help gives the rule with a physical equivalent, which is the number that lets you compare it against the ergonomics literature in A6:

> "Any on-screen element that someone can click, touch, or otherwise interact with should be large enough for reliable interaction. Consider making sure these elements have a width and height of at least 48dp, as described in the Material Design Accessibility guidelines."

> "Ensure that each of those elements is 48x48dp in size, or approximately 9mm in each dimension."

Source: Google, *Touch target size* (Android Accessibility Help), <https://support.google.com/accessibility/android/answer/7101858>.

Google also makes the point that the target is the *hit area*, not the drawn area:

> "Touch targets include the area that responds to user input. Touch targets extend beyond the visual bounds of an element: An element like an icon may appear to be 24x24dp but the padding surrounding it comprises the full 48x48dp touch target."

Same source.

**Unit note.** 1 dp equals 1 CSS px at a 160 dpi baseline, so 48dp reads as 48 CSS px in a web view.

**Not found.** Material Design 3 publishes no point-of-sale page. Material's density guidance is aimed at desktop data tables, not at touch order entry.

### A3. WCAG 2.2

**Class: hard standard.** This is the only source in Part A that is normative and legally referenced. Text below is quoted verbatim from the W3C Recommendation.

#### Target size

**SC 2.5.8 Target Size (Minimum), Level AA:**

> "The size of the target for pointer inputs is at least 24 by 24 CSS pixels, except when:
> - **Spacing** Undersized targets (those less than 24 by 24 CSS pixels) are positioned so that if a 24 CSS pixel diameter circle is centered on the bounding box of each, the circles do not intersect another target or the circle for another undersized target;
> - **Equivalent** The function can be achieved through a different control on the same page that meets this criterion;
> - **Inline** The target is in a sentence or its size is otherwise constrained by the line-height of non-target text;
> - **User Agent Control** The size of the target is determined by the user agent and is not modified by the author;
> - **Essential** A particular presentation of the target is essential or is legally required for the information being conveyed."

**SC 2.5.5 Target Size (Enhanced), Level AAA:**

> "The size of the target for pointer inputs is at least 44 by 44 CSS pixels except when: Equivalent ... Inline ... User Agent Control ... Essential ..."

Source: W3C, *Web Content Accessibility Guidelines (WCAG) 2.2*, <https://www.w3.org/TR/WCAG22/>. Understanding document for 2.5.8: <https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html>.

Read the 24 CSS px figure as a **legal floor, not a design target**. The Spacing exception means a 20 px control can still pass 2.5.8 if nothing else is near it. Neither reading is a good POS button.

#### Contrast

**SC 1.4.3 Contrast (Minimum), Level AA:**

> "The visual presentation of text and images of text has a contrast ratio of at least 4.5:1, except for the following: Large Text Large-scale text and images of large-scale text have a contrast ratio of at least 3:1 ..."

**SC 1.4.6 Contrast (Enhanced), Level AAA:**

> "The visual presentation of text and images of text has a contrast ratio of at least 7:1, except for the following: Large Text Large-scale text and images of large-scale text have a contrast ratio of at least 4.5:1 ..."

**SC 1.4.11 Non-text Contrast, Level AA:**

> "The visual presentation of the following have a contrast ratio of at least 3:1 against adjacent color(s): User Interface Components Visual information required to identify user interface components and states ... Graphical Objects Parts of graphics required to understand the content ..."

"Large scale (text)" is defined normatively as:

> "with at least 18 point or 14 point bold or font size that would yield equivalent size for Chinese, Japanese and Korean (CJK) fonts"

Source, all four: <https://www.w3.org/TR/WCAG22/>.

**Important limit on 1.4.3 for this project.** The 4.5:1 figure was derived for a display viewed in normal indoor light. WCAG says nothing about ambient light washing out the rendered contrast. See A5.

### A4. Vendor POS design documentation

**Class: vendor guidance. The headline result is how little of it exists.**

Six POS vendors were searched for public design documentation with hard numbers on button size, grid density, category colour coding, or order-entry flow. Two publish a usable number. Four publish none.

| Vendor | Public design guidance? | Hard numbers |
| --- | --- | --- |
| Square | No design guidance | none |
| Toast | Nothing found | none |
| Clover | Yes, qualitative | one contrast ratio, borrowed from WCAG |
| Shopify POS | Yes, qualitative | none |
| Lightspeed | Nothing public | none |
| Oracle MICROS Simphony | Yes | capacity limits only |

**Square.** The Point of Sale API documentation covers only launching the Square app and reading the callback. It is integration documentation, not design documentation. <https://developer.squareup.com/docs/pos-api/how-it-works>. A document titled "Square POS Design System" circulates online; it is a designer's personal portfolio piece, not a Square publication, and must not be cited as vendor guidance.

**Toast.** Nothing found. Searches for Toast POS design guidance return only unrelated design-system pages about the "toast" notification component. Toast publishes no discoverable public UI design guidance.

**Clover.** <https://docs.clover.com/dev/docs/app-design-requirements>. Guidance is qualitative, with one number, and Clover then defers outright:

> "a contrast ratio of at least 4.5:1"

That figure is WCAG SC 1.4.3, not a Clover-original number. Clover also states that it "recommends referring to Google Material Design guidelines as a good place to start". Clover does give per-device orientation rules: Station is landscape, Flex is portrait.

**Shopify POS.** <https://shopify.dev/docs/api/pos-ui-extensions>. Defers to Polaris web components for visual consistency. No pixel or dp numbers of its own.

**Lightspeed.** No public numeric UI design guidance. Lightspeed has an internal design system called Flame, documented only on a former employee's portfolio site. The O-Series "design your POS look and layout" help article is merchant-facing, about choosing colours and logos, not a developer spec.

**Oracle MICROS Simphony.** The one vendor that publishes numbers relevant to grid density, and they are **capacity limits, not ergonomic recommendations**:

> "Approximately 50 tabs and subtabs can appear on a page"

> "each tab and subtab can contain up to 50 buttons"

Source: Oracle, *Simphony Configuration Guide - Workstation Touchscreen Page Design*, <https://docs.oracle.com/en/industries/food-beverage/simphony/19.3/simcg/c_workstation_touchscreen_page_design.htm>.

Read these as what the software will tolerate, not as what an operator can use. Fifty buttons on one page is a database limit. Nothing in Oracle's documentation says fifty buttons is a good idea.

**Conclusion for lead A4.** There is no such thing as published POS touch-UI design guidance with defensible numbers. The POS industry does not publish design specifications. A developer who wants numbers must take them from the platform guidelines (A1, A2), the accessibility standard (A3), and the ergonomics literature (A6).

### A5. Outdoor and high-ambient-light legibility

**Class: hard standard (partly paywalled) plus weak vendor material.**

This is the weakest-sourced part of Part A, and the report says so plainly.

**ISO 9241-303:2011, *Ergonomics of human-system interaction - Part 303: Requirements for electronic visual displays*.** The full standard is paywalled. The publicly readable preview at <https://cdn.standards.iteh.ai/samples/57992/bddfd91165b444f6b9815a6993feadc5/ISO-9241-303-2011.pdf> carries the front matter and the early clauses; the later clauses on luminance contrast and contrast for object legibility appear in the table of contents but the preview truncates before their text. Catalogue entry: <https://www.iso.org/standard/57992.html>.

What the readable part does give, quoted verbatim:

> "In the ambient illumination for which the display is designed, the display luminance shall exceed the minimum value for obtaining a sufficient recognizability of the displayed information over the design viewing range and the intended lifetime of the visual display unit."

> "For an office application having 500 lx illuminance (horizontally) of white paper with a reflectance of 80 % and positive display polarity, it is often recommended that the display luminance be in the range of 100 cd/m2 to 150 cd/m2."

> "For an intended uniform display luminance, the luminance non-uniformity, either step-wise or smooth, in ambient illumination shall not exceed the threshold for reduced visual performance, with a maximum of 1,7:1."

> "The area average luminance of task areas that are frequently viewed in sequence while using the display (paper document, screen, etc.) should be between 0,1L and 10L, where L is the average luminance of the display."

And, relevant to a tablet mounted on a stand rather than held:

> "a display shall be legible from any angle of inclination up to at least 40 degrees from the normal to the surface of the display."

That clause matters here because off-axis contrast on an LCD collapses fastest exactly where reflected sunlight is worst. A tablet on a fixed stand, viewed by an Operator of unpredictable height, is the case the clause is written for.

The load-bearing point for an outdoor booth is the first quote. ISO ties the required display luminance to the ambient illuminance the display is *designed for*. Direct sunlight is roughly 100 000 lx against the 500 lx office case in the second quote, that is two to three orders of magnitude more light. A display specified for an office does not meet ISO 9241-303 outdoors.

**The one peer-reviewed source found.** Zhang, K., Han, T., Cho, W.-K., Kwok, H.-S., and Liu, Z. (2021). "Investigation of Enhanced Ambient Contrast Ratio in Novel Micro/Mini-LED Displays." *Nanomaterials*, 11(12), 3304. DOI [10.3390/nano11123304](https://doi.org/10.3390/nano11123304). Reported to state that LED panels "with a brightness of over 7000 nits ... can have adequate readability outdoor under a clear sky", while standard LCDs in the hundreds of nits read adequately only below roughly 500 lux. **Flagged as not independently verified**: MDPI returned HTTP 403 to automated retrieval, so this quote comes from a research assistant's reading of the paper, not from a fetch performed here. The DOI is real and the paper is open access; verify in a browser before relying on the figure.

Note how far that 7000 nits figure is from a consumer tablet, which is typically 500 to 600 nits. **No tablet on the market is a sunlight-readable display.** Shade at the booth is a more effective intervention than any UI decision in this report.

**What could not be sourced.** No free, primary, citable source was found that states a required display luminance in cd/m2 for outdoor readability, nor a required effective contrast ratio under a stated ambient illuminance. The figures that circulate widely - "over 1000 nits", "contrast ratio 800:1", "5:1 or higher in high ambient light" - trace only to display-vendor marketing pages, for example <https://www.displaymodule.com/blogs/knowledge/sunlight-readable-tft-displays-nits-brightness-contrast-anti-glare>. **Treat these as vendor marketing, not as standards.** They are recorded here only so that a later reader does not mistake them for evidence.

**Physics that does not need a citation, and is safe to design against.** Ambient light reflected off the screen adds a roughly constant luminance to both the light and the dark parts of the image. Adding a constant to both terms of a contrast ratio always moves the ratio towards 1:1. So the contrast a designer measures in a browser is an upper bound, and the contrast an Operator sees in daylight is always lower. The practical consequence is in the Implications section: design well above the WCAG floor, and prefer dark text on a light ground, because the light ground is the term that ambient reflection hurts least in relative terms.

### A6. Fitts's law and empirical touch-target studies

**Class: academic finding.** This is the most valuable part of Part A, because it is the only body of evidence produced by measuring people rather than by committee.

The literature splits sharply by **task posture and touch strategy**, and that split is the key to reading it.

- **Lift-off strategy, one-handed thumb, small phone.** Small targets work. Roughly 9 to 10 mm.
- **Land-on strategy, standing, finger, kiosk.** Small targets do not work. Roughly 20 mm.

An Apple Fest POS is the second case, not the first.

#### The most directly relevant study

**Colle, H. A., and Hiszem, K. J. (2004).** "Standing at a kiosk: effects of key size and spacing on touch screen numeric keypad performance and user preference." *Ergonomics*, 47(13), 1406-1423. DOI [10.1080/00140130410001724228](https://doi.org/10.1080/00140130410001724228). PubMed: <https://pubmed.ncbi.nlm.nih.gov/15513716/>.

Quoted verbatim from the abstract, retrieved through the NCBI E-utilities API:

> "Twenty participants used finger touches to enter one, four or 10 digits in a numeric keypad displayed on a capacitive touch screen, while standing in front of a touch screen kiosk. Key size (10, 15, 20, 25 mm square) and edge-to-edge key spacing (1, 3 mm) were factorially combined. ... **Spacing had no measurable effects.** Entry times were longer and errors were higher for smaller key sizes, but no significant differences were found between key sizes of 20 and 25 mm. Participants also preferred 20 mm keys to smaller keys, and they were indifferent between 20 and 25 mm keys. **Therefore, a key size of 20 mm was found to be sufficiently large for land-on key entry.**"

Two findings matter here, and the second is the surprising one.

1. **20 mm square is the number**, and going past it buys nothing.
2. **Spacing had no measurable effect** across 1 mm and 3 mm gaps. Size did the work, not the gap. This directly contradicts the platform guidance in A1 and A2, which treats padding as roughly co-equal with size. The reconciliation is that Colle and Hiszem only tested gaps up to 3 mm, and a 3 mm gap is far below the 12 pt / 8dp padding those guidelines recommend. So the honest reading is: **within a small range, size dominates spacing.** It is not evidence that spacing never matters.

#### Supporting studies at the kiosk end

- **Sears, A., and Shneiderman, B. (1991).** "High precision touchscreens: design strategies and comparisons with a mouse." *International Journal of Man-Machine Studies*, 34(4), 593-613. DOI [10.1016/0020-7373(91)90037-8](https://doi.org/10.1016/0020-7373%2891%2990037-8). With a lift-off strategy, targets as small as 1.7 x 2.2 mm rivalled a mouse; but the fastest and most accurate condition overall was the largest target tested, 13.8 x 17.9 mm. With a land-on strategy, targets need to be about 20 mm square.
- **Duff, S., Irwin, C., Skye, J., Sesto, M., and Wiegmann, D. (2010).** "The Effect of Disability and Approach on Touch Screen Performance During a Number Entry Task." *Proceedings of the Human Factors and Ergonomics Society 54th Annual Meeting.* DOI [10.1037/e578672012-006](https://doi.org/10.1037/e578672012-006). Tested 10, 20, 25 and 30 mm. Concluded buttons should be at least 20 mm, with minimal gain beyond that.
- **Chen, K. B., et al. (2013).** "Touch screen performance by individuals with and without motor control disabilities." *Applied Ergonomics*, 44(2), 297-302. DOI [10.1016/j.apergo.2012.08.004](https://doi.org/10.1016/j.apergo.2012.08.004). Tested 10, 20, 30 mm with 1 and 3 mm gaps. Performance for users without motor impairment plateaued at 20 mm; for users with motor impairment it kept improving past 20 mm.
- **Tao, D., Yuan, J., Liu, S., and Qu, X. (2018).** "Effects of button design characteristics on performance and perceptions of touchscreen use." *International Journal of Industrial Ergonomics.* DOI [10.1016/j.ergon.2017.12.001](https://doi.org/10.1016/j.ergon.2017.12.001). Tested nine button sizes from 10 mm to 50 mm and four spacing intervals on a 65-inch touchscreen, with an older and a younger age group. Larger buttons reduced response time for both groups.

Four independent studies converge on 20 mm. That convergence is the strongest single result in this whole report.

**Provenance.** The Colle and Hiszem abstract above was fetched and quoted directly from the NCBI E-utilities API. The bibliographic details for Sears and Shneiderman, Duff et al., and Chen et al. were confirmed against the Crossref API. Their reported findings, however, come from a research assistant's reading and were **not** re-verified against the full texts, which are paywalled. Treat the 20 mm convergence as strong and the individual per-study details as good but second-hand.

#### The mobile-thumb end, for contrast

- **Parhi, P., Karlson, A. K., and Bederson, B. B. (2006).** "Target size study for one-handed thumb use on small touchscreen devices." *Proceedings of MobileHCI 2006*, 203-210. DOI [10.1145/1152215.1152260](https://doi.org/10.1145/1152215.1152260). Recommends 9.2 mm for discrete tasks and 9.6 mm for serial tasks, with no significant error-rate gain beyond those.
- **Bi, X., Li, Y., and Zhai, S. (2013).** "FFitts law: modeling finger touch with Fitts' law." *Proceedings of CHI 2013*, 1363-1372. DOI [10.1145/2470654.2466180](https://doi.org/10.1145/2470654.2466180). Models finger-touch endpoints as a dual distribution and reports that its index of difficulty explains more variance than classic Fitts's law for finger input (R-squared at or above 0.91).

Note the coincidence and do not be misled by it: Apple's 44 pt is about 8.5 mm on an 11-inch iPad and Material's 48dp is 7.62 mm nominal, so both land in the **7 to 9 mm** range. That is the *thumb-on-a-phone* number, not the *finger-at-a-kiosk* number. The platform guidelines were written for phones. A festival till is a kiosk.

#### Not found

No published study applies Fitts's law to a full POS checkout screen, that is, to a grid of menu-item buttons with an order pane. Every touch-target study found uses a numeric keypad or an abstract target-acquisition task. Anyone claiming a "Fitts-optimal POS layout" from published work is extrapolating.

---

## PART B - the young-operator question

**Short answer: almost nothing exists, and the gap is real rather than a search failure.**

No standards body, no platform owner, no POS vendor, and no HCI publication issues user-interface design guidance targeted at operators aged 14 to 16. What does exist falls into three groups, none of which is design guidance:

1. **Operating guidance embedded in youth fundraising apps.** Girl Scouts and Scouts BSA both ship apps that 14- to 16-year-olds use to take payments. Both carry task guidance. Neither carries any statement that the interface was designed differently because the operator is young.
2. **Regulation that permits the work and says nothing about the tool.** See B6.
3. **A touchscreen motor-performance literature with a documented hole at exactly this age band.** See B5. This is the finding with the most weight.

### B1. Girl Scouts, Digital Cookie

**Class: organisational operating guidance. The strongest source in Part B, and it is still not design guidance.**

Girl Scout grade levels map to age: Cadettes are grades 6 to 8, **Seniors are grades 9 to 10, roughly 14 to 16**, Ambassadors are grades 11 to 12. Source: <https://www.girlscouts.org/en/discover/about-us/what-girl-scouts-do/grade-levels.html>. So Seniors match the target band almost exactly, and they use the same Digital Cookie app as everyone else. **There is no age-differentiated interface.**

The two official GSUSA training PDFs are:

- *Using the Mobile App: A Training Guide for Families at a Cookie Booth*, <https://www.girlscouts.org/content/dam/girlscouts-gsusa/forms-and-documents/cookie/digitalcookie/Mobile_App_Booth.pdf>
- *Using the Digital Cookie Mobile App: A Training Guide for Caregivers and Girl Scouts*, <https://www.girlscouts.org/content/dam/girlscouts-gsusa/forms-and-documents/cookie/digitalcookie/MobileAppCaregivers.pdf>

Quoted operating guidance that is genuinely aimed at a young booth operator:

> "Tip: Have good lighting and double check the numbers before placing the order." (**verified**, appears in both PDFs)

> "Check with your troop cookie volunteer before approving any orders through this feature" (**verified**)

> "For booth sales, the delivery type must be 'Give now.'" (**verified**)

Two further quotes were reported by a research assistant and could **not** be verified here, because they appear inside screenshot images in the PDFs rather than in extractable text. They are recorded with that flag:

> "Do not use public wi-fi to send your order" (**unverified**)

> "Do not hand your mobile device to the customer." (**unverified**)

The three verified quotes were extracted directly from the PDFs with `pdftotext`.

Three of these are design decisions in disguise, and they are worth copying:

- **The device stays with the operator.** No hand-over-the-tablet flow.
- **Higher-stakes actions route through an adult.** The routine cash sale is left to the young operator alone; approval of anything unusual is gated.
- **The interface constrains the order type at a booth** rather than trusting the operator to pick correctly.

What is absent: no target size, no contrast figure, no statement of motor-accuracy reasoning, no document anywhere saying "this screen looks like this because our users are 14".

### B2. Scouts BSA and the Trail's End app

**Class: organisational operating guidance. Weak.**

Scouts BSA covers ages 11 to 17 (<https://www.scouting.org/programs/scouts-bsa/>), so the band is in scope. The Trail's End help article *How to Enter a Storefront Sale*, <https://support.trails-end.com/en/articles/13401129-how-to-enter-a-storefront-sale>, describes the flow: "Make a Sale", then "Continue with Storefront Sale", add products, pay by card, customer device, or cash, then an optional receipt. One operating rule is stated:

> "Both cash and credit card sales should be entered in real time,"

No age-specific guidance of any kind. The article does not even state a minimum age for independent operation. **Confirms the age group uses such tools; contributes nothing on design.**

### B3. DECA school-based enterprise

**Class: organisational certification standard. Weak and tangential.**

The DECA School-Based Enterprise Certification Program Guidelines, <https://dpi.wi.gov/sites/default/files/imce/mmee/pdf/deca_sbe_certification_guidelines.pdf>, contain two POS-adjacent performance indicators:

> "Open/Close register/terminal." and "Prepare cash drawers/banks." (Standard 1, Financial Analysis)

> "Are instructions for equipment (food heating stations, POS systems, other examples) conveniently displayed? Are SBE employees trained on proper use of equipment?" (Standard 2, Operations, mandatory indicator)

DECA therefore acknowledges that high-school students operate POS systems, and requires that instructions be "conveniently displayed". That is an accreditation requirement about process, not an interface specification. It says nothing about button size, workflow, or error tolerance. **No national standard for school stores with POS specifications was found.**

### B4. 4-H, FFA, Junior Achievement

**Class: absence of evidence.**

- **4-H**: a "Business and Citizenship" curriculum collection exists (<https://shop4-h.org/collections/business-citizenship-curriculum>) but contains no retail or POS operations material.
- **FFA**: nothing found.
- **Junior Achievement**: JA BizTown (<https://jausa.ja.org/programs/ja-biztown>) has children role-play storefronts, but it is elementary-age, below the target band. JA Company Program is high-school age but publishes only general business-venture material.

No published document from any of the three addresses point-of-sale interface design or operation for youth of any age.

### B5. Academic literature on adolescent touchscreen performance

**Class: confirmed absence of evidence. This is the most important result in Part B.**

Touchscreen and pointing performance research clusters into two disconnected bands with a hole between them:

- **Early childhood, ages roughly 3 to 10.** Well studied.
- **Adults, roughly 18 upward through the elderly.** Very well studied.
- **Ages 11 to 17.** Effectively unstudied.

Evidence for each side of the hole:

**The childhood side.**
- Vatavu, R.-D., Cramariuc, G., and Schipor, D. M. (2015). "Touch interaction for children aged 3 to 6 years: Experimental findings and relationship to motor skills." *International Journal of Human-Computer Studies*, 74, 54-76. 89 children aged 3 to 6, plus 30 young adults for comparison. <https://mintviz.usv.ro/publications/ijhcs2015.pdf>
- Yadav, et al. (2021). "Children's interaction with touchscreen devices: Performance and validity of Fitts' law." *Human Behavior and Emerging Technologies*, 3(5), 1132-1140. 30 children aged 4 to 10; concludes that their interaction "does not obey Fitts' law".

**The clearest single piece of evidence that the gap is real.** The MTAGIC project (Lisa Anthony, University of Florida, NSF-funded, six studies) is the most comprehensive children's-touchscreen research programme in HCI. Its own summary states:

> "Initial studies were broad in their target age range (ages 5-17), but later studies focused on the age range of 5-10 years old."

Source: <https://mtagic.wordpress.com/>. **Flagged as not independently verified**: this sentence was reported by a research assistant but was not found in the extractable text of the project home page, so it may sit on a sub-page or in a linked paper. The published record is consistent with it either way - see "Physical dimensions of children's touchscreen interactions: Lessons from five years of study on the MTAGIC project", *International Journal of Human-Computer Studies*, <https://www.sciencedirect.com/science/article/abs/pii/S1071581918302441>, whose studies are on ages 5 to 10.

If the quote holds, a project *scoped to reach 17-year-olds* narrowed to 5 to 10 and published no adolescent data. **The wider gap finding does not depend on this quote**; it rests on the absence of any located study covering the band.

**The nearest miss on the adult side.** Hertzum, M., and Hornbaek, K. (2010). "How Age Affects Pointing with Mouse and Touchpad: A Comparison of Young, Adult, and Elderly Users." *International Journal of Human-Computer Interaction*, 26(7), 703-734. <https://mortenhertzum.dk/publ/IJHCI2010b.pdf>. Its three groups are:

> "young (12-14 years), adult (25-33 years), and elderly (61-69 years)"

(**Verified**: extracted directly from the PDF.)

This is the only located study whose "young" bracket touches the band at all, and it (a) stops at 14, and (b) tests **mouse and touchpad, not touchscreen**. The authors state their own purpose in terms that make the general problem explicit:

> "this study provides for direct comparison of young and elderly users and for more precisely assessing the biases introduced by basing design on adult users only - the predominant user group in studies of pointing devices and input techniques."

One large real-world logging study, "Temporal clusters of age-related behavioral alterations captured in smartphone touchscreen interactions", <https://pmc.ncbi.nlm.nih.gov/articles/PMC9418599/>, has N=598 across ages 16 to 86, but the distribution is bimodal with modes at 25 and 63, and no result is reported separately for the youngest participants.

**Statement of the gap, plainly.** No touchscreen study was found with participants aged 14 to 16 as a distinct, analysed group. Any claim that 14- to 16-year-olds need larger, smaller, or otherwise different touch targets than adults is **unsupported by published evidence in either direction**. The honest design position is that this age group should be treated as adults for motor performance, and as novices for everything else.


### B6. US Department of Labor youth employment rules

**Class: hard regulation.** This is the only part of Part B with a firm, quotable answer, and the answer is the opposite of what the question implies.

**Cashiering is expressly permitted for 14- and 15-year-olds.** 29 CFR 570.34 lists the occupations that minors aged 14 and 15 may perform. Paragraph (d) reads, verbatim and in full:

> "(d) Cashiering, selling, modeling, art work, work in advertising departments, window trimming, and comparative shopping."

Source: 29 CFR 570.34, "Occupations that may be performed by minors 14 and 15 years of age", <https://www.ecfr.gov/current/title-29/subtitle-B/chapter-V/subchapter-A/part-570/subpart-C/section-570.34>. Retrieved 2026-09-01 through the eCFR renderer API, because the eCFR web page redirects an automated fetcher.

The same section opens by making the permission conditional on hours and on the prohibited-occupation lists:

> "This subpart authorizes only the following occupations in which the employment of minors 14 and 15 years of age is permitted when performed for periods and under conditions authorized by 570.35 and not involving occupations prohibited by 570.33 ..."

So the constraint on a 14- or 15-year-old at a till is a **time-of-day and hours-per-week constraint** (29 CFR 570.35), not a task constraint. Nothing in the regulation restricts what a cash register may show them or what it may let them do.

**And the regulation probably does not apply here at all.** The FLSA reaches employees. A Scout working a troop booth at a festival is a volunteer. The Department of Labor states:

> "Individuals who volunteer or donate their services, usually on a part-time basis, for public service, religious or humanitarian objectives, not as employees and without contemplation of pay, are not considered employees ..."

Source: US Department of Labor, elaws FLSA Advisor, *Volunteers*, <https://webapps.dol.gov/elaws/whd/flsa/docs/volunteers.asp>.

One caveat that a reader should not skip. DOL Fact Sheet #14A, *Non-Profit Organizations and the Fair Labor Standards Act*, <https://www.dol.gov/agencies/whd/fact-sheets/14a-flsa-non-profits>, states, verbatim:

> "The FLSA recognizes the generosity and public benefits of volunteering and allows individuals to freely volunteer in many circumstances for charitable and public purposes. Individuals may volunteer time to religious, charitable, civic, humanitarian, or similar non-profit organizations as a public service and not be covered by the FLSA. **Individuals generally may not, however, volunteer in commercial activities run by a non-profit organization such as a gift shop.** A volunteer generally will not be considered an employee for FLSA purposes if the individual volunteers freely for public service, religious or humanitarian objectives, and without contemplation or receipt of compensation. Typically, such volunteers serve on a part-time basis and do not displace regular employed workers or perform work that would otherwise be performed by regular employees."

(dol.gov returns HTTP 403 to automated fetchers. The text above was retrieved verbatim from the Internet Archive capture: <https://web.archive.org/web/2024id_/https://www.dol.gov/agencies/whd/fact-sheets/14a-flsa-non-profits>.)

A festival food booth that sells to the public bears a resemblance to the gift-shop example. Whether the FLSA reaches it turns on whether the sale is the nonprofit's ordinary commercial activity or incidental fundraising. **No DOL document addressing a one-time community festival booth was found**, so this is an inference from adjacent guidance, not a ruling. It is also not a design question: even in the worst case, the constraint is hours worked, not screen layout.

The mirror list is equally clear. 29 CFR 570.33, the prohibited occupations for 14- and 15-year-olds (<https://www.law.cornell.edu/cfr/text/29/570.33>), contains **no mention of cash registers, POS terminals, or checkstands**. The prohibitions cover manufacturing, mining, power-driven machinery, motor vehicles, ladders, baking, freezer and meat-cooler work, and similar physical hazards.

Hours limits for 14- and 15-year-olds come from 29 CFR 570.35, summarised in DOL Fact Sheet #43, <https://www.dol.gov/agencies/whd/fact-sheets/43-child-labor-non-agriculture>: outside school hours only, not before 7 a.m. or after 7 p.m., extended to 9 p.m. from 1 June through Labor Day, up to 8 hours on a non-school day and 40 hours in a non-school week. **Flagged as not independently fetched** - dol.gov blocks automated retrieval and this summary came from search extracts of that page. Note that Apple Fest runs roughly 8 a.m. to 6 p.m., inside the 7 a.m. to 9 p.m. summer window.

**Bottom line for the design.** No US federal rule limits what a POS may ask a 14- to 16-year-old to do. The regulation is not a design input. Design difficulty, not legality, is the real constraint.


### B7. HCI literature on novice operators and error recovery

**Class: academic. Thin, and adult-focused.**

No study of **novice cashiers as operators** was found. The nearest analogue in the literature is the self-checkout customer, who is an untrained person driving a till, but that is a different task and a different population.

- "An analysis into early customer experiences of self-service checkouts: Lessons for improved usability." A diary study with 31 respondents unfamiliar with self-checkouts, documenting difficulty, cognitive strain, and errors. <https://www.researchgate.net/publication/332641235_An_analysis_into_early_customer_experiences_of_self-service_checkouts_Lessons_for_improved_usability>
- "The Impact of Guidance Information on Checkout Efficiency and User Experience in Supermarket Self-Checkout Machines" (2025). Found photo-based guidance cues outperformed illustrated cues. <https://journals.sagepub.com/doi/10.1177/21582440251404195>
- **Cognitive walkthrough** (Lewis, Polson, Wharton, and Rieman, 1990) is the evaluation method built for exactly this problem. It was created to evaluate walk-up-and-use interfaces such as kiosks and ATMs, where a user's ability to succeed with no prior knowledge and no training is what matters. <https://www.nngroup.com/articles/cognitive-walkthroughs/>. It is a **method, not a finding**, and it carries no age-specific results, but it is the right method to apply to this POS.

None of this disaggregates by age. Commercial "cashier training" material from POS vendors is practitioner blog content with no method and no data; it is not cited here.

### Absence of evidence: what was searched

Recorded so that a later reader does not repeat the search.

| Target | Where searched | Result |
| --- | --- | --- |
| Age-specific POS or till UI guidance, 14-16 | Apple HIG, Material 3, WCAG, Square, Toast, Clover, Shopify, Lightspeed, Oracle | Nothing. No vendor or platform mentions operator age at all. |
| Youth-operator design guidance | girlscouts.org, digitalcookie.girlscouts.org, trails-end.com, scouting.org | Task and safety guidance only. No design rationale tied to age. |
| Youth retail programme POS specs | deca.org, dpi.wi.gov DECA guidelines, shop4-h.org, ffa.org, jausa.ja.org | DECA mentions POS systems in a process requirement. Nothing else. |
| Touchscreen motor accuracy, ages 14-16 | web search, whose results surfaced ACM DL, ScienceDirect, PubMed and Google Scholar; terms: children touch target size, age differences touchscreen target acquisition, adolescent fine motor touchscreen, Fitts law children touchscreen, teenagers touchscreen performance | **No study with 14-16 as a distinct analysed group.** Literature clusters at ages 3-10 and 18-86. |
| Novice cashier error rates | web search, whose results surfaced ACM DL, SAGE and ResearchGate; terms: novice cashier, POS usability study, cash register usability error, training-free interface, walk-up-and-use, error recovery novice users | No study of employed novice cashiers. Only self-checkout customer studies. |
| Fitts's law applied to a POS menu grid | web search, ACM DL and Google Scholar results | Not found. All touch-target studies use numeric keypads or abstract targets. |
| Outdoor display luminance requirement, citable | ISO catalogue, ANSI/HFES, Google Scholar | ISO 9241-303 is paywalled past the preview. Free numeric figures are display-vendor marketing. |


---

## Contradictions between sources

These sources do not agree. Where they conflict, the conflict is recorded rather than averaged away.

**1. Apple contradicts Apple.** The Buttons page says a button "needs a hit region of at least 44x44 pt". The Accessibility page says 44x44 pt is the *default* for iOS and iPadOS and the *minimum* is 28x28 pt. Both are current. Use 44 pt as the floor.

**2. Every source gives a different number, and the units hide it.** Converted to physical size (see the unit note below):

| Source | Stated | Physical, approximately |
| --- | --- | --- |
| WCAG 2.2 SC 2.5.8 (AA) | 24 x 24 CSS px | 4.6 mm at 132 px/in |
| Material 3 / Android | 48 x 48 dp | 7.6 mm nominal; Google rounds it to "approximately 9mm" |
| Apple iOS / iPadOS | 44 x 44 pt | 8.5 mm on a ~11-inch iPad |
| WCAG 2.2 SC 2.5.5 (AAA) | 44 x 44 CSS px | 8.5 mm at 132 px/in |
| Parhi et al. 2006, thumb on a phone | 9.2 to 9.6 mm | 9.2 to 9.6 mm |
| **Colle and Hiszem 2004, finger at a standing kiosk** | **20 mm** | **20 mm** |

Note what this table shows. **Apple, Google, and WCAG AAA all cluster around 8 to 9 mm, which is the mobile-phone thumb number.** The kiosk literature says roughly double that. The platform guidelines were written for a phone in one hand. They are not wrong; they answer a different question. **A festival till is a standing kiosk, so 20 mm governs.**

Google's own "approximately 9mm" gloss on 48dp is itself generous. 48dp at the 160 dpi baseline is 7.62 mm.

**3. Colle and Hiszem contradict Apple and Material on spacing.** Colle and Hiszem found "Spacing had no measurable effects" across 1 mm and 3 mm gaps. Apple recommends about 12 pt of padding; Material recommends 8dp of separation. The reconciliation: Colle and Hiszem only tested gaps up to 3 mm, well below what the platform guidelines recommend, and their key sizes were large. The defensible reading is **size dominates spacing within a narrow range**, not that spacing is free. Keep a real gutter; it costs little and it satisfies the WCAG 2.5.8 Spacing exception without effort.

**4. WCAG's own floor is the lowest number anyone gives, on purpose.** 24 CSS px is a conformance floor across all content types, not a recommendation for a task-critical touch UI. Its Spacing exception even permits smaller targets. Do not design to it.

**5. Clover defers to Google; Oracle publishes limits, not advice.** Clover's only number, 4.5:1, is WCAG's. Oracle's numbers, 50 tabs and 50 buttons per tab, are database capacity limits that no one should mistake for a density recommendation.

**6. ISO 9241-303 is scoped to the seated office.** Its worked example assumes 500 lx and a 300 mm-plus viewing distance. An outdoor booth is two orders of magnitude brighter. The standard's requirement still applies in principle - luminance must suit the design ambient illuminance - but its numbers do not.

---

## A note on units, because the numbers above are not comparable as written

Points, dp, CSS pixels, and millimetres are four different things, and the conversion depends on the device.

- **1 pt (iOS/iPadOS) = 1 CSS px** in a web view at default zoom.
- **1 dp (Android) = 1 CSS px** at the 160 dpi baseline.
- **CSS px to millimetres depends on the display.** The CSS reference density is 96 px per inch, but a real tablet is denser. A ~10.9-inch iPad presents about 1180 x 820 CSS px on a 10.9-inch diagonal, which is about **132 CSS px per inch**.

The conversion is:

```
css_px = mm * (css_px_per_inch / 25.4)
```

Worked both ways for the governing 20 mm figure:

| Device CSS density | 20 mm becomes |
| --- | --- |
| 96 CSS px/in (CSS reference) | 76 CSS px |
| 132 CSS px/in (~11-inch iPad) | **104 CSS px** |

A developer can measure their own device's density in the browser:

```js
const cssPxPerInch = Math.hypot(screen.width, screen.height) / DIAGONAL_INCHES;
const px = mm => Math.round(mm * cssPxPerInch / 25.4);
```

**Do not use the CSS `mm` unit for this.** CSS `mm` is defined against the 96 px/in reference pixel, not against the physical display, so `20mm` in CSS renders as 76 CSS px regardless of the actual panel. It will be too small on a tablet.


---

## Implications for the Apple Fest POS

The recommendations below are numeric and directly applicable. Each names the source it comes from. Where sources conflict, the strictest applicable one is chosen, and the reason is given.

**Two facts drive all of it.**

1. The Apple Fest till is a **standing kiosk operated with a finger**, not a phone held in one hand. The kiosk literature governs, not the platform guidelines. (A6)
2. The operator being 14 to 16 changes **nothing about target sizing**, because no evidence exists in either direction (B5). It changes a great deal about **error recovery and cognitive load**, because the operator is a novice working a queue. Design for the novice, not for the age.

### Touch targets

| Property | Recommendation | Source and reason |
| --- | --- | --- |
| Primary item button, physical | **20 mm x 20 mm minimum** | Colle and Hiszem 2004 (20 mm sufficient for land-on entry at a standing kiosk); corroborated by Duff et al. 2010 and Chen et al. 2013. Four studies converge. |
| Primary item button, CSS px on a ~11-inch iPad | **104 x 104 CSS px minimum** | 20 mm at ~132 CSS px/in. Recompute for any other device with the formula above. |
| Absolute floor for any control | **88 x 88 CSS px** (about 17 mm) | 15 mm was measurably worse than 20 mm in Colle and Hiszem; 17 mm is an **interpolated** floor, not a measured one. |
| Secondary controls (quantity plus/minus, category tabs) | **still 88 CSS px minimum**, never smaller | These are the most-tapped controls in a queue. Do not shrink them to save space. |
| Destructive controls (void line, void order) | **104 CSS px, and separated by at least 48 CSS px from any non-destructive control** | Beyond published guidance. Rationale: the WCAG Spacing exception logic applied to consequence rather than to size. |
| Do not use | 24 CSS px (WCAG floor), 44 pt (Apple), 48 dp (Material) | These are phone-thumb figures or conformance floors. They are all under half the kiosk requirement. (Contradiction 2) |

### Spacing

| Property | Recommendation | Source |
| --- | --- | --- |
| Gutter between adjacent item buttons | **16 CSS px minimum** | Apple's 12 pt padding guidance, rounded up; Material's 8dp is the weaker figure. Colle and Hiszem found spacing did not matter in the 1 to 3 mm range, so this is cheap insurance, not a measured requirement. |
| Padding inside a button, around its label | **12 CSS px minimum** | Apple HIG, "about 12 points of padding around elements that include a bezel". |
| Gap between the item grid and the order pane | **32 CSS px minimum** | Prevents a mis-aimed item tap landing in the order list. |

### Contrast, sized for daylight rather than for an office

| Property | Recommendation | Source and reason |
| --- | --- | --- |
| Button label against button fill | **7:1 minimum** | WCAG SC 1.4.6 (AAA), not the 4.5:1 AA floor. Reason: ambient light reflected off the panel adds a constant to both terms of the contrast ratio, which always pushes the measured ratio towards 1:1. The 4.5:1 figure assumes indoor viewing. Design with headroom. (A3, A5) |
| Body and secondary text | **7:1 minimum** | Same. |
| Button borders, focus rings, selected states | **4.5:1 minimum** | WCAG SC 1.4.11 requires 3:1. Raised for the same daylight reason. |
| Polarity | **Dark text on a light fill** | Positive polarity. Reflected ambient light degrades a light ground proportionally less than it degrades a dark one. ISO 9241-303's own worked example assumes positive polarity. |
| Colour as an information channel | **Never alone** | Category colour coding is fine as reinforcement. Every category must also carry a text label. WCAG SC 1.4.1. Sunlight and glare distort perceived hue. |
| Button label size | **24 CSS px minimum, 600 weight or heavier** | 24 CSS px is 18 pt, which qualifies as "large scale (text)" under WCAG. It also keeps a label legible at arm's length on a stand. |

### Grid density and layout

Assume a landscape tablet at roughly 1180 x 820 CSS px.

| Property | Recommendation | Reason |
| --- | --- | --- |
| Item button pitch | **120 CSS px** (104 button + 16 gutter) | Derived from the two tables above. |
| Maximum columns in the item grid | **6**, given about 800 CSS px of grid width after the order pane | 800 / 120 = 6.67. |
| Maximum rows visible | **5**, given about 700 CSS px of grid height after header and footer | 700 / 120 = 5.8. |
| Practical ceiling on visible items | **30**, and aim for **12 to 16** | 6 x 5 is the geometric limit. A festival menu is small. Fit it in one screen. |
| Scrolling in the item grid | **None.** If the menu does not fit, cut the menu or use category tabs, never a scroll. | A scrolling grid makes button positions non-constant, which destroys the muscle memory a novice operator builds during a rush. |
| Category tabs | **At most 4**, always visible, never nested | Oracle's 50-tabs figure is a capacity limit, not advice. (A4) |
| Button positions | **Fixed for the whole event day.** No reordering by popularity, no "recently used" reshuffling. | Same muscle-memory reason. This is the single highest-value rule for a novice operator and it costs nothing. |

### Novice-operator rules, adapted from the youth-sales apps

These come from B1 and B7. They are the closest thing to age-relevant guidance that exists, and none of them is a size number.

1. **The device never leaves the operator's hands.** Girl Scouts state this outright: "Do not hand your mobile device to the customer." Do not build a customer-facing confirmation step that requires passing the tablet.
2. **Route the unusual through an adult.** Girl Scouts gate order approval behind "Check with your troop cookie volunteer". Apply the same shape here: a young Operator takes orders freely; voiding a completed order, opening the till, or reprinting should be a Leader action or at least require a deliberate second step.
3. **Constrain the flow so the wrong choice is not offered.** Girl Scouts force "Give now" at a booth rather than trusting the operator to pick the delivery type. Do the same: at the booth there is one order type, and the interface should not present others.
4. **Every action must be undoable in one tap, and undo must be as large as the action.** No published number supports a size here; the reason is that a novice under queue pressure makes errors of commission, and the recovery path must not itself require precision.
5. **Confirm with a number, not with a dialogue.** Girl Scouts' guidance is "double check the numbers before placing the order". Show the running total large and constantly, rather than asking "Are you sure?" at the end.
6. **Evaluate the screens with a cognitive walkthrough**, the method built for walk-up-and-use interfaces where the user has had no training (B7). This is the correct evaluation method for this project and it needs no participants.

### What this document does not license anyone to claim

- That 14- to 16-year-olds need bigger buttons than adults. **No evidence exists.** (B5)
- That any law restricts what the POS may let a young operator do. **Cashiering is expressly permitted, and volunteers are outside the FLSA anyway.** (B6)
- That there is an industry-standard POS layout. **The POS industry publishes no design specifications.** (A4)
- That "1000 nits" or any similar figure is a standard for outdoor readability. **It is display-vendor marketing.** (A5)
