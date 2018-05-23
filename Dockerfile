FROM scratch
ADD post-scrapper /
ENTRYPOINT ["/post-scrapper"]