# `txn-cli`

> Convert bank statements to YNAB-compatible CSV

`txn-cli` is a command-line tool to work with bank statements in CSV format.

The purpose of this tool is focused on _conversion_ of one CSV format to
another. It is NOT to validate bookkeeping. That is why this project is built
with the philosophy of lenient-but-loud, meaning that it will generally allow
invalid, nonsensical, or even corrupted data, but will try it's best to warn you
or give you hints about problems it finds.
