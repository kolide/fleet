if ! git diff-files --quiet ./docs
   then
   echo 'Dependencies documentation must be updated. Please run `make lint-license` locally and commit documentation changes.'
   exit 1
fi
