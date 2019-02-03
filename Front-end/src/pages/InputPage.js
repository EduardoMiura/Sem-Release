import React from "react";
import {
  MDBInput,
  MDBInputSelect,
  MDBFormInline,
  MDBBtn,
  MDBContainer,
  MDBRow,
  MDBCol
} from "mdbreact";
import DocsLink from "./DocsLink";

class InputPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      value: "John Doe"
    };
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleSubmit(event) {
    alert("MDBInput value: " + this.state.value);
    const data = new FormData(event.target);

    fetch('/api/form-submit-url', {
      method: 'POST',
      body: data,
    });
    event.preventDefault();
  }

  saveToState = value => {
    this.setState({ ...this.state, value: value });
  };

  handleChange = value => {
    console.log(value);
  };

  render() {
    return (
      <MDBContainer className="mt-5">
        <DocsLink
          title="Inputs"
          href="https://mdbootstrap.com/docs/react/forms/inputs/"
        />
        <MDBContainer style={{ textAlign: "initial" }}>
          <div>

            <h4 className="mt-4">
              <strong>Horizontal form</strong>
            </h4>
            <form >
              <div className="form-group row">
                <label
                  htmlFor="inputEmail3"
                  className="col-sm-2 col-form-label"
                >
                  InstallationID
                </label>
                <div className="col-sm-10">
                  <input
                    type="number"
                    className="form-control"
                    id="inputEmail3"
                    placeholder="InstallationID"
                  />
                </div>
              </div>
              <div className="form-group row">
                <label
                  htmlFor="inputPassword3"
                  className="col-sm-2 col-form-label"
                >
                  Integration
                </label>
                <div className="col-sm-10">
                  <input
                    type="number"
                    className="form-control"
                    id="inputPassword3"
                    placeholder="IntegrationID"
                  />
                </div>
              </div>
              <div className="form-group row">
                <label
                  htmlFor="filePem"
                  className="col-sm-2 col-form-label"
                >
                  Arquivo PEM
                </label>
                <div className="col-sm-10">
                  <input
                    id="file" type="file" />
                </div>
              </div>


              <div className="form-group row">
                <div className="col-sm-10">
                  <button type="submit" className="btn btn-primary btn-md">
                    Enviar
                  </button>
                </div>
              </div>
            </form>

          </div>
        </MDBContainer>
      </MDBContainer>
    );
  }
}

export default InputPage;
