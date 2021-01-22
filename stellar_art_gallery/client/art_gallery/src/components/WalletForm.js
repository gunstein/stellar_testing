import React, { useState} from "react";
import Button from "@material-ui/core/Button";
import TextField from "@material-ui/core/TextField";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import { createAccount } from "../services/stellar";

export default function WalletForm({ handleCloseToParent, paymentinfo }) {
  const [secretkey, setSecretkey] = useState("");
  const [publickey, setPublickey] = useState("");
  const [payMode, setPayMode] = useState(false);

  const handleClose = () => {
    handleCloseToParent();
  };

  const handleCreateAccount = () => {
    createAccount().then(account_response =>{
        setSecretkey(account_response.secret)
        setPublickey(account_response.publickey)
        setPayMode(true)
    })
  };

  const handlePay = () => {
    handleCloseToParent();
  };

  return (
    <Dialog
      open={true}
      onClose={handleClose}
      aria-labelledby="form-dialog-title"
    >
      <DialogTitle id="form-dialog-title">
        Pay
      </DialogTitle>
      <DialogContent>
        <DialogContentText>
          Make account and pay. (This dialog is only relevant for the stellar test network.)
        </DialogContentText>
            <TextField
                id="target-account-pay"
                label="Target account"
                value={targetAccount}
                fullWidth={true}
                InputProps={{
                    readOnly: true,
                }}
            />
            <TextField
                id="target-memo-pay"
                label="Target memo"
                value={targetMemo}
                InputProps={{
                    readOnly: true,
                }}
            />        
            <TextField
                id="account-public-key"
                label="Account public key"
                value={publickey}
                fullWidth={true}
                InputProps={{
                    readOnly: true,
                }}
            />
            <TextField
                id="account-secret-key"
                label="Account secret key"
                value={secretkey}
                fullWidth={true}
                InputProps={{
                    readOnly: true,
                }}
            />  
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          Cancel
        </Button>
        {!payMode ? (
          <Button onClick={handleCreateAccount} color="primary">
            Create account
          </Button>
        ) : null}  
        {payMode ? (
          <Button onClick={handlePay} color="primary">
            Pay
          </Button>      
        ) : null}        
      </DialogActions>
    </Dialog>
  );
}
