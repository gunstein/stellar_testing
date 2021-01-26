import StellarSdk from "stellar-sdk";

export async function createAccount() {
  const pair = StellarSdk.Keypair.random();
  try {
    const response = await fetch(
      `https://friendbot.stellar.org?addr=${encodeURIComponent(
        pair.publicKey()
      )}`
    );
    const responseJSON = await response.json();
    console.log("SUCCESS! You have a new account :)\n", responseJSON);
    return pair;
  } catch (e) {
    console.error("ERROR!", e);
  }
}

export async function getAccountBalance(account) {
  const server = new StellarSdk.Server("https://horizon-testnet.stellar.org");
  const account_response = await server.loadAccount(account);
  return { balance: account_response.balances[0].balance };
}

export async function sendPayment(sourceKeyPair, targetAccount, memo, amount) {
  console.log("sendPayment: targetAccount : ", targetAccount);
  console.log("sendPayment: memo : ", memo);
  console.log("sendPayment: amount : ", amount);
  //Mostly copied from here https://developers.stellar.org/docs/tutorials/send-and-receive-payments/
  var StellarSdk = require("stellar-sdk");
  var server = new StellarSdk.Server("https://horizon-testnet.stellar.org");

  //var destinationId = "GA2C5RFPE6GCKMY3US5PAB6UZLKIGSPIUKSLRB6Q723BM2OARMDUYEJ5";
  // Transaction will hold a built transaction we can resubmit if the result is unknown.
  var transaction;

  // First, check to make sure that the destination account exists.
  // You could skip this, but if the account does not exist, you will be charged
  // the transaction fee when the transaction fails.
  server
    .loadAccount(targetAccount)
    // If the account is not found, surface a nicer error message for logging.
    .catch(function (error) {
      if (error instanceof StellarSdk.NotFoundError) {
        throw new Error("The destination account does not exist!");
      } else return error;
    })
    // If there was no error, load up-to-date information on your account.
    .then(function () {
      return server.loadAccount(sourceKeyPair.publicKey());
    })
    .then(function (sourceAccount) {
      // Start building the transaction.
      transaction = new StellarSdk.TransactionBuilder(sourceAccount, {
        fee: StellarSdk.BASE_FEE,
        networkPassphrase: StellarSdk.Networks.TESTNET,
      })
        .addOperation(
          StellarSdk.Operation.payment({
            destination: targetAccount,
            // Because Stellar allows transaction in many currencies, you must
            // specify the asset type. The special "native" asset represents Lumens.
            asset: StellarSdk.Asset.native(),
            amount: amount.toString(10),
            //amount: "1",
          })
        )
        // A memo allows you to add your own metadata to a transaction. It's
        // optional and does not affect how Stellar treats the transaction.
        .addMemo(StellarSdk.Memo.text(memo))
        // Wait a maximum of three minutes for the transaction
        .setTimeout(180)
        .build();
      // Sign the transaction to prove you are actually the person sending it.
      transaction.sign(sourceKeyPair);
      // And finally, send it off to Stellar!
      return server.submitTransaction(transaction);
    })
    .then(function (result) {
      console.log("Success! Results:", result);
      return result;
    })
    .catch(function (error) {
      console.error("Something went wrong!", error);
      // If the result is unknown (no response body, timeout etc.) we simply resubmit
      // already built transaction:
      // server.submitTransaction(transaction);
    });
}
