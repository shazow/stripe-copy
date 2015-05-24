# stripe-copy

**Status**: *Unstable. Use with caution, please audit the code and contribute.*

Sometimes you need to migrate between two different Stripe accounts. You can ask Stripe to copy your customer objects, but they will not copy the rest of the objects. See: [How can I migrate to a new Stripe account?](https://support.stripe.com/questions/how-can-i-migrate-to-a-new-stripe-account-7a206563-51ad-4c70-a995-a01f57a3eb56)

`stripe-copy` is a command-line tool for copying Stripe objects like Plans and Subscriptions between accounts.


## Usage

The tool loads your private keys from environment variables `STRIPE_FROM` and `STRIPE_TO`.

```shell
export STRIPE_FROM="YOUR_PRIVATE_API_KEY" STRIPE_TO="OTHER_PRIVATE_API_KEY"
```

```shell
stripe-copy --help
```


## License

MIT.
