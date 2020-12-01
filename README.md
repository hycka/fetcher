# Fetcher

Fetch news from urls.  

# Intro

Fetch depth: 1  

# Development Tips

There files need to be modified while add a new site:
- main.go: add entrance url
- links.go -> SetLinks(): add case about target urls feature regex, eg: if url must have `about`, param 2 is `.*?about.*`
- site/newsite: copy files from sibling folder, then develop and pass the test.
- post.go -> TreatPost(): add case for new site domain.
