import React from "react";

import "./HomePage.css";
//import REPOSITORY_METHODS from '../../assets/scripts/api';

import {
  Container,
  Col,
  Row,
  Fa
} from "mdbreact";
const NavLink = require("react-router-dom").NavLink;

// function fecthRepository() {
//   const resposeOfAPI = REPOSITORY_METHODS.getRepositories();
//   return handleResponse(resposeOfAPI, getMesssage, dispatch);
// }

// function handleResponse(response, errorMessage) {
//   return response
//     .then((valueOfResponse) => {
//       if (getStatus(valueOfResponse.status)) {
//         return Promise.reject(new Error('Promise error on handleResponse'));
//       }
//       return valueOfResponse.data;
//     })
//     .catch(
//       error => Promise.reject(errorMessage(error))
//       ,
//     );
// }
class HomePage extends React.Component {
  render() {
    //console.log('SSSSSSS', fecthRepository())
    return (
      <div>

        <Container>
          <Row>
            <Col md="10" className="mx-auto mt-4">

              <div>
                <NavLink
                  tag="button"
                  className="btn btn-sm indigo darken-3"
                  to="/forms/input"
                >
                  Adicionar Nova Configuração
                    </NavLink>
              </div>
              <hr className="my-5" />
              <Row id="categories">
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="css3" className="pink-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">CSS</h4>
                    <p className="grey-text">
                      Animations, colours, shadows, skins and many more! Get to
                      know all our css styles in one place.
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo darken-3"
                      to="/css"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="cubes" className="blue-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">COMPONENTS</h4>
                    <p className="grey-text">
                      Ready-to-use components that you can use in your
                      applications. Both basic and extended versions!
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo lighten-2"
                      to="/components"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="code" className="green-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">ADVANCED</h4>
                    <p className="grey-text">
                      Advanced components such as charts, carousels, tooltips
                      and popovers. All in Material Design version.
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo darken-3"
                      to="/advanced"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
              </Row>
              <Row id="categories">
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="bars" className="pink-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">NAVIGATION</h4>
                    <p className="grey-text">
                      Ready-to-use navigation layouts, navbars, breadcrumbs and
                      much more! Learn more about our navigation components.
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo darken-3"
                      to="/navigation"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="edit" className="blue-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">FORMS</h4>
                    <p className="grey-text">
                      Inputs, autocomplete, selecst, date and time pickers.
                      Everything in one place is ready to use!
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo lighten-2"
                      to="/forms"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="table" className="green-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">TABLES</h4>
                    <p className="grey-text">
                      Basic and advanced tables. Responsive, datatables, with
                      sorting, searching and export to csv.
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo darken-3"
                      to="/tables"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
              </Row>
              <Row id="categories">
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="window-restore" className="pink-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">MODALS</h4>
                    <p className="grey-text">
                      Modals used to display advanced messages to the user.
                      Cookies, logging in, registration and much more.
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo darken-3"
                      to="/modals"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="arrows" className="blue-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">EXTENDED</h4>
                    <p className="grey-text">
                      Google Maps, Social Buttons, Pre-built Contact Forms and
                      Steppers. Find out more about our extended components.
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo lighten-2"
                      to="/extended"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
                <Col md="4" className="mb-5">
                  <Col size="2" md="2" className="float-left">
                    <Fa icon="th" className="green-text" size="2x" />
                  </Col>
                  <Col size="10" md="8" lg="10" className="float-right">
                    <h4 className="font-weight-bold">SECTIONS</h4>
                    <p className="grey-text">
                      E-commerce, contact, blog and much more sections. All
                      ready to use in seconds.
                    </p>
                    <NavLink
                      tag="button"
                      className="btn btn-sm indigo darken-3"
                      to="/sections"
                    >
                      Learn more
                    </NavLink>
                  </Col>
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </div>
    );
  }
}

export default HomePage;
