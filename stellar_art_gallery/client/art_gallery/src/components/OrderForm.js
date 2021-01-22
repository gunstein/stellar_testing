import React, { useState} from "react";
import Button from "@material-ui/core/Button";
import TextField from "@material-ui/core/TextField";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import { createOrder } from "../services/order";
import {WalletForm} from "./WalletForm"

export default function OrderForm({ handleCloseToParent, tile }) {
  //const [open, setOpen] = React.useState(false);
  const [emailInput, setEmailInput] = useState("");
  const [payMode, setPayMode] = useState(false);
  const [walletMode, setWalletMode] = useState(false);
  const [targetAccount, setTargetAccount] = useState("");
  const [targetMemo, setTargetMemo] = useState("");

  const handleClose = () => {
    handleCloseToParent();
  };

  const handleBuy = () => {
    createOrder({email:emailInput, artid:tile.artid})
    .then(order_response => {
      console.log("order_response: ", order_response)
      if(order_response.hasOwnProperty("data") && order_response.data.hasOwnProperty("account") && order_response.data.hasOwnProperty("memo"))
      {
        setTargetAccount(order_response.data.account)
        setTargetMemo(order_response.data.memo)
        setPayMode(true)
      }
    })
    //handleCloseToParent();
  };

  const handlePay = () => {
    setWalletMode(true)
    //handleCloseToParent();
  };

  return (
    <Dialog
      open={true}
      onClose={handleClose}
      aria-labelledby="form-dialog-title"
    >
      <DialogTitle id="form-dialog-title">
        Buy {tile.title} by {tile.artist}
      </DialogTitle>
      <DialogContent>
        <DialogContentText>
          Please enter email and press buy. Then paymentinfo will show up. Download info will be sent per email when payment is received.
        </DialogContentText>
        {!payMode ? (
          <TextField
            autoFocus
            margin="dense"
            id="email"
            label="Email Address"
            type="email"
            fullWidth
            onChange={event => setEmailInput(event.target.value)} 
            value={emailInput}
          />
        ) : null}
        {payMode ? (
          <React.Fragment>
            <TextField
              id="target-account"
              label="Target account"
              value={targetAccount}
              fullWidth={true}
              InputProps={{
                readOnly: true,
              }}
            />
            <TextField
              id="target-memo"
              label="Target memo"
              value={targetMemo}
              InputProps={{
                readOnly: true,
              }}
            />  
          </React.Fragment>        
        ) : null}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          Cancel
        </Button>
        {!payMode ? (
          <Button onClick={handleBuy} color="primary">
            Buy
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
