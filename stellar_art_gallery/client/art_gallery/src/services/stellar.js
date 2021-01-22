import StellarSdk from "stellar-sdk";


export function createAccount() {
    const pair = StellarSdk.Keypair.random();
    return fetch(`https://friendbot.stellar.org?addr=${encodeURIComponent(
        this.pair.publicKey()
        )}`).then(function(data){
            console.log(data)
            return pair.json() //only interested in the pair
        }).catch(function(error) {
            console.log("Failed!", error)
        })   
}

