# spacer

a spaced repetition program

## TODOS

- [ ] Problem: let the user quickly say whether they got the answer right or not. 
    * Solution: after reveal, show a section that says: I got it... and then a green button "right" and a red button "wrong"
- [ ] Problem: sort flashcards by least well remembered 
    * Solution: record a count of correct answers with each flashcard. when user answers "right", increment by one. when user answers "wrong", reset to zero. sort by the correctness count, ascending. record a date of last answered, and sort by that, as a secondary sort, ascending.
- [ ] Problem: give the user some sense of progress / accomplishment by indicating that they don't need to redo cards that they've completed successfully.
    * Solution: record a date of last successfully reviewed. add a ready-for-review section at the top. If there are none, put a "you're all set 🎉" message there. add a next-flashcards section. If there are none, don't show this section. If there are no cards at all, don't show either section. Cards go in the ready-for-review when their check-in date is today or before. Cards have a check-in date set to today when they are created or when they are answered incorrectly. Cards have a check-in date set according to the following formula: fib[correct answers]+today.
- [ ] Problem: resurface difficult cards more frequently.
    * Solution: record a wrong-answer count. update the target date formula to be fib[correct answers]/wrong answers + today. increment wrong-answer count whenever there's a wrong answer.
- [ ] Problem: provide a way for a user to just leave themselves a note. 
    * Solution: update the q/a prompts to be q/subject, a/note.
- [ ] Problem: provide a way to edit a flashcard 
    * Solution: add an edit button with a pencil icon. replace the q/a with editable form fields, and an "update" button. Reset "correct answers" to zero.
- [ ] Problem: provide a way to set a reminder/todo. 
    * Solution: add an optional target date field when editing a card. add a subnote describing that it will remind you more frequently the closer you are to the target date. check-in date formula is fib counting back from the target date. instead of remembered/not, an ack button. When you ack on or after the target date, shift the target date to tomorrow.
- [ ] Problem: provide sync between platforms.
    * Solution: add a google drive integration? add a folder for cards. title each card as the creation date/time. toml as storage format.
- [ ] Problem: not available to people. 
    * Solution: deploy a public-facing website.
- [ ] Problem: no way to organize. 
    * Solution: allow tags. display all in-use tags in two sections under the add card interface: included/not-included. display all tags for a card at the bottom of that card. when adding or editing a card, clicking a tag will move it from included/not-included. there should be an edit field for adding a new tag. existing tags are filtered down while typing. entering an existing tag should just reuse the existing tag rather than creating a new one. When typing in any of the add fields, filter the list to matches
- [x] Problem: this doesn't drag & drop on mobile
    * Solution: ??? ask copilot to fix that? yep. prompt was "The flashcards are not draggable on mobile. Make the flashcards draggable on mobile."