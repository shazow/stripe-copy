# stripe-copy

**Status**: *Unstable. Use with caution, please audit the code and contribute.*

Sometimes you need to migrate between two different Stripe accounts. You can ask Stripe to copy your customer objects, but they will not copy the rest of the objects. See: [How can I migrate to a new Stripe account?](https://support.stripe.com/questions/how-can-i-migrate-to-a-new-stripe-account-7a206563-51ad-4c70-a995-a01f57a3eb56)

> We are able to do one-time copies of customer data across accounts, so that you don’t have to ask your users to re-enter their credit card details in the event that you are migrating to a new account. Get in touch to arrange copying this data over. [...] Charges, invoices, plans, subscriptions, coupons, events, and logs do not get copied over. Only the raw customer objects get copied over.

`stripe-copy` is a command-line tool for copying Stripe objects like Plans and Subscriptions between accounts.


## Install

You'll need Go to build the source, ideally version 1.4 or newer.

```
$ go get github.com/shazow/stripe-copy
```

If this project matures, we'll add binaries with tagged releases.


## Usage

The tool loads your private keys from environment variables `STRIPE_SOURCE` and
`STRIPE_TARGET`.

```
$ export STRIPE_SOURCE="YOUR_PRIVATE_API_KEY" STRIPE_TARGET="OTHER_PRIVATE_API_KEY"
```

```
$ stripe-copy --help
Usage:
  stripe-copy [OPTIONS]

Application Options:
  -v, --verbose  Show verbose logging.
  -p, --pretend  Do everything read-only, skip writes.
      --version

Help Options:
  -h, --help     Show this help message
```

```
$ stripe-copy -vv --pretend
2015-05-29 12:50:45.870 INFO Running in pretend mode. Write operations will be skipped.
2015-05-29 12:50:45.870 DEBUG Loading target plans...
2015-05-29 12:50:47.119 DEBUG Loaded 5 target plans. Loading source plans...
2015-05-29 12:50:47.338 INFO Plans: 5 loaded, 0 missing, 0 changed.
2015-05-29 12:50:47.338 DEBUG Pretend mode: Stopping early.
```

## Roadmap

In approximate order of priority:

- [x] Sync Plans
- [x] Confirm Customers
- [ ] Sync Subscriptions (if customers are present)
- [ ] Release v1.0
- [ ] Optionl to cancel subscriptions on source
- [ ] Export to file (YAML?)
- [ ] Import from file
- [ ] Release v1.1
- [ ] Parallelize
- [ ] Setup continuous builds for binaries of tagged releases
- [ ] Refactor sync logic into a library and document
- [ ] Add tests
- [ ] Option to delete target items missing from source
- [ ] Option to do a bi-directional sync

Anything else we'd like to sync?


## Contributing

Please do. Even if you're just getting started with Go, we're happy to guide you
along.

1. **[Check for open issues](https://github.com/shazow/stripe-copy/issues) or open
   a fresh issue** to start a discussion around a feature idea or a bug. Let us know
   what you'll be working to avoid duplicating effort. Bonus points if you pick a task off the roadmap above.
2. **Send a pull request.** It's okay to do this early to get some feedback. Add
   `[WIP]` to the title if you're not done yet.
3. Politely yet persistently bug the maintainer until it gets reviewed and merged. :)

**Too little time but plenty of money?** [Add a bounty on
Bountysource](https://www.bountysource.com/).

## License

MIT.
