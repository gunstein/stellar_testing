import React, { useState} from "react";
import Button from "@material-ui/core/Button";
import TextField from "@material-ui/core/TextField";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import { createOrder } from "../services/order";

export default function FormDialog({ handleCloseToParent, tile }) {
  //const [open, setOpen] = React.useState(false);
  const [emailInput, setEmailInput] = useState('');

  const handleClose = () => {
    handleCloseToParent();
  };

  const handleBuy = () => {
    createOrder({email:emailInput, artid:tile.artid})
    .then(order_response => {
        console.log("order_response: ", order_response) 
    })
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
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label="Email Address"
          type="email"
          fullWidth
          onChange={event => setEmailInput(event.target.value)} 
          value={emailInput}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          Cancel
        </Button>
        <Button onClick={handleBuy} color="primary">
          Buy
        </Button>
      </DialogActions>
    </Dialog>
  );
}
