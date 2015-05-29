# stripe-copy

**Status**: *Unstable. Use with caution, please audit the code and contribute.*

Sometimes you need to migrate between two different Stripe accounts. You can ask Stripe to copy your customer objects, but they will not copy the rest of the objects. See: [How can I migrate to a new Stripe account?](https://support.stripe.com/questions/how-can-i-migrate-to-a-new-stripe-account-7a206563-51ad-4c70-a995-a01f57a3eb56)

> We are able to do one-time copies of customer data across accounts, so that you don’t have to ask your users to re-enter their credit card details in the event that you are migrating to a new account. Get in touch to arrange copying this data over. [...] Charges, invoices, plans, subscriptions, coupons, events, and logs do not get copied over. Only the raw customer objects get copied over.

`stripe-copy` is a command-line tool for copying Stripe objects like Plans and Subscriptions between accounts.


## Usage

The tool loads your private keys from environment variables `STRIPE_SOURCE` and
`STRIPE_TARGET`.

```shell
$ export STRIPE_SOURCE="YOUR_PRIVATE_API_KEY" STRIPE_TARGET="OTHER_PRIVATE_API_KEY"
```

```shell
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

```shell
$ stripe-copy -vv --pretend
2015-05-29 12:50:45.870 INFO Running in pretend mode. Write operations will be skipped.
2015-05-29 12:50:45.870 DEBUG Loading target plans...
2015-05-29 12:50:47.119 DEBUG Loaded 5 target plans. Loading source plans...
2015-05-29 12:50:47.338 INFO Plans: 5 loaded, 0 missing, 0 changed.
2015-05-29 12:50:47.338 DEBUG Pretend mode: Stopping early.
```

## Roadmap

- [x] Sync Plans
- [ ] Confirm Customers
- [ ] Sync Subscriptions (if customers are present)
- [ ] Export to file (YAML?)
- [ ] Import from file
- [ ] Option to delete target items missing from source
- [ ] Option to do a bi-directional sync

Anything else we'd like to sync? Open an issue or send a pull request.


## License

MIT.
