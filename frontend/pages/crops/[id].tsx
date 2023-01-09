import { useRef, useState } from "react";
import type { NextPage } from "next";
import Link from "next/link";
import {
  Button,
  Card,
  Col,
  Form,
  InputGroup,
  ListGroup,
  Row,
  Tab,
  Table,
  Tabs,
} from "react-bootstrap";
import {
  FaLongArrowAltLeft,
  FaPlus,
  FaUtensilSpoon,
  FaTrash,
  FaPaperPlane,
  FaCamera,
  FaCut,
  FaExchangeAlt,
} from "react-icons/fa";

import ButtonIcon from "../../components/ButtonIcon";
import ModalContainer from "../../components/ModalContainer";
import Panel from "../../components/Panel";
import TableTaskItem from "../../components/TableTaskItem";
import Layout from "../../components/Layout";
import { tasksData, notesData } from "../../data";
import useModal from "../../hooks/useModal";

const CropDetail: NextPage = () => {
  const { modalOpen, showModal, closeModal } = useModal();
  const [dueDate, setDueDate] = useState("");
  const [priority, setPriority] = useState("");
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("");
  const [isError, setIsError] = useState(false);
  const target = useRef(null);

  const addTask = () => {
    if (!dueDate || !priority || !title) {
      setIsError(true);
    } else {
      setIsError(false);
      closeModal();
    }
  };

  // Harvest Modal
  const [modalHarvestOpen, setModalHarvestOpen] = useState(false);
  const [harvestArea, setHarvestArea] = useState("");
  const [harvestType, setHarvestType] = useState("");
  const [harvestQty, setHarvestQty] = useState("");
  const [harvestUnit, setHarvestUnit] = useState("");
  const [harvestNotes, setHarvestNotes] = useState("");

  // Move Modal
  const [modalMoveOpen, setModalMoveOpen] = useState(false);
  const [moveSource, setMoveSource] = useState("");
  const [moveDest, setMoveDest] = useState("");
  const [moveQty, setMoveQty] = useState(0);

  // Dump Modal
  const [modalDumpOpen, setModalDumpOpen] = useState(false);
  const [dumpArea, setDumpArea] = useState("");
  const [dumpQty, setDumpQty] = useState(0);
  const [dumpNotes, setDumpNotes] = useState("");

  // Take Picture Modal
  const [modalPicOpen, setModalPicOpen] = useState(false);
  const [pic, setPic] = useState("");
  const [picNotes, setPicNotes] = useState("");

  return (
    <Layout>
      <Row>
        <Col className="mb-3">
          <Link href="/crops">
            <ButtonIcon
              label="Back to Crops Batches"
              icon={<FaLongArrowAltLeft className="me-2" />}
              onClick={() => {}}
              variant="link"
            />
          </Link>
        </Col>
      </Row>
      <Row>
        <Col>
          <Tabs defaultActiveKey="basic">
            <Tab eventKey="basic" title="Basic Info">
              <h3 className="py-3">Romaine</h3>
              <Row>
                <Col className="mb-3">
                  <small>Batch ID</small>
                  <div>
                    <strong>rom-24mar</strong>
                  </div>
                </Col>
                <Col className="mb-3">
                  <small>Initial Planning</small>
                  <div>
                    <strong>777 Post on Organic lettuce</strong>
                  </div>
                </Col>
              </Row>
              <Row>
                <Col className="mb-3">
                  <small>Status</small>
                  <div>
                    <strong>0 Seeding, 777 Growing, 0 Dumped</strong>
                  </div>
                </Col>
                <Col className="mb-3">
                  <small>Current Quantity</small>
                  <div>
                    <strong>777 Post on Organic lettuce</strong>
                  </div>
                </Col>
              </Row>
              <Row>
                <Col className="mb-3">
                  <small>Seeding Date</small>
                  <div>
                    <strong>24/03/2021</strong>
                  </div>
                </Col>
              </Row>
              <Row>
                <Col className="mb-3">
                  <small>Last Watering</small>
                  <div>
                    <strong>-</strong>
                  </div>
                </Col>
              </Row>
              <Row>
                <Col className="mb-3">
                  <div className="d-grid gap-2">
                    <ButtonIcon
                      label="Harvest"
                      icon={<FaCut className="me-2" />}
                      onClick={() => setModalHarvestOpen(true)}
                      variant="secondary"
                      isBlock
                    />
                  </div>
                </Col>
                <Col className="mb-3">
                  <div className="d-grid gap-2">
                    <ButtonIcon
                      label="Move"
                      icon={<FaExchangeAlt className="me-2" />}
                      onClick={() => setModalMoveOpen(true)}
                      variant="primary"
                      isBlock
                    />
                  </div>
                </Col>
                <Col className="mb-3">
                  <div className="d-grid gap-2">
                    <ButtonIcon
                      label="Dump"
                      icon={<FaTrash className="me-2" />}
                      onClick={() => setModalDumpOpen(true)}
                      variant="danger"
                      isBlock
                    />
                  </div>
                </Col>
                <Col className="mb-3">
                  <div className="d-grid gap-2">
                    <ButtonIcon
                      label="Take Picture"
                      icon={<FaCamera className="me-2" />}
                      onClick={() => setModalPicOpen(true)}
                      variant="light"
                      isBlock
                    />
                  </div>
                </Col>
              </Row>
              <h5 className="py-3">Activity</h5>
              <Card>
                <Card.Body>
                  <div className="d-flex">
                    <div>
                      <FaUtensilSpoon className="me-2" />
                    </div>
                    <div>
                      <div>
                        Seeded <strong>777 Pots</strong> of rom-24mar on{" "}
                        <strong>Organic lettuce</strong>
                      </div>
                      <small className="mt-1 text-muted">
                        24/03/2021 at 14:51
                      </small>
                    </div>
                  </div>
                </Card.Body>
              </Card>
            </Tab>
            <Tab eventKey="notes" title="Tasks &amp; Notes">
              <Row>
                <Panel title="Tasks">
                  <>
                    <div className="mb-3">
                      <ButtonIcon
                        label="Add Task"
                        icon={<FaPlus className="me-2" />}
                        onClick={() => showModal()}
                        variant="primary"
                      />
                    </div>
                    <Table responsive>
                      <thead>
                        <tr>
                          <th />
                          <th className="w-75">Items</th>
                          <th>Category</th>
                        </tr>
                      </thead>
                      <tbody>
                        {tasksData &&
                          tasksData.map(
                            ({
                              id,
                              item,
                              details,
                              dueDate,
                              priority,
                              category,
                            }) => (
                              <tr key={id}>
                                <td>
                                  <Form>
                                    <Form.Check type="checkbox" />
                                  </Form>
                                </td>
                                <td>
                                  <TableTaskItem
                                    id={id}
                                    item={item}
                                    details={details}
                                    dueDate={dueDate}
                                    priority={priority}
                                  />
                                </td>
                                <td>
                                  <span className="text-uppercase">
                                    {category}
                                  </span>
                                </td>
                              </tr>
                            )
                          )}
                      </tbody>
                    </Table>
                  </>
                </Panel>
                <Panel title="Notes">
                  <>
                    <InputGroup className="mb-3">
                      <Form.Control type="text" placeholder="Create a note" />
                      <Button variant="secondary">
                        <div className="d-flex align-items-center">
                          <FaPaperPlane />
                        </div>
                      </Button>
                    </InputGroup>
                    <ListGroup>
                      {notesData &&
                        notesData.map(({ id, title, createdOn }) => (
                          <ListGroup.Item key={id}>
                            <div className="d-flex align-items-center justify-content-between py-1">
                              <div>
                                <div className="mb-1">{title}</div>
                                <small className="text-muted">
                                  {createdOn}
                                </small>
                              </div>
                              <div>
                                <FaTrash />
                              </div>
                            </div>
                          </ListGroup.Item>
                        ))}
                    </ListGroup>
                  </>
                </Panel>
              </Row>
            </Tab>
          </Tabs>
        </Col>
      </Row>
      <ModalContainer
        title="Crop: Add New Task from rom-24mar"
        isShow={modalOpen}
        handleCloseModal={closeModal}
        handleSubmitModal={addTask}
      >
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Due Date</Form.Label>
            <InputGroup ref={target}>
              <Form.Control
                type="date"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
              />
            </InputGroup>
            {isError && (
              <Form.Text className="text-danger">
                The due date field is required
              </Form.Text>
            )}
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Is this task urgent?</Form.Label>
            <Form.Check
              type="radio"
              label="Yes"
              name="priority"
              onChange={() => setPriority("urgent")}
            />
            <Form.Check
              type="radio"
              label="No"
              name="priority"
              onChange={() => setPriority("normal")}
            />
            {isError && (
              <Form.Text className="text-danger">
                The priority field is required
              </Form.Text>
            )}
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Task Category</Form.Label>
            <Form.Select onChange={(e) => setSelectedCategory(e.target.value)}>
              <option>Please select category</option>
              <option value="1">Reservoir</option>
              <option value="2">Pest Control</option>
              <option value="3">Safety</option>
              <option value="4">Sanitation</option>
            </Form.Select>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Title</Form.Label>
            <Form.Control
              type="text"
              onChange={(e) => setTitle(e.target.value)}
              value={title}
            />
            {isError && (
              <Form.Text className="text-danger">
                The title field is required
              </Form.Text>
            )}
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Description</Form.Label>
            <Form.Control
              as="textarea"
              onChange={(e) => setDesc(e.target.value)}
              value={desc}
              style={{ height: "120px" }}
            />
          </Form.Group>
        </Form>
      </ModalContainer>

      <ModalContainer
        title="Harvest rom-24mar"
        isShow={modalHarvestOpen}
        handleCloseModal={() => setModalHarvestOpen(false)}
        handleSubmitModal={() => setModalHarvestOpen(false)}
      >
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Choose area to be harvested</Form.Label>
            <Form.Select onChange={(e) => setHarvestArea(e.target.value)}>
              <option>Please select area</option>
              <option value="1">Organic lettuce</option>
              <option value="2">Organic chilli</option>
            </Form.Select>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Choose type of harvesting</Form.Label>
            <Form.Select onChange={(e) => setHarvestType(e.target.value)}>
              <option>All</option>
              <option value="1">Partial</option>
            </Form.Select>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Quantity</Form.Label>
            <Form.Control
              type="input"
              value={harvestQty}
              onChange={(e) => setHarvestQty(e.target.value)}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Units</Form.Label>
            <Form.Select onChange={(e) => setHarvestUnit(e.target.value)}>
              <option value="1">Grams</option>
              <option value="2">Kilograms</option>
            </Form.Select>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Notes</Form.Label>
            <Form.Control
              as="textarea"
              onChange={(e) => setHarvestNotes(e.target.value)}
              value={harvestNotes}
              style={{ height: "120px" }}
            />
          </Form.Group>
        </Form>
      </ModalContainer>

      <ModalContainer
        title="Move rom-24mar"
        isShow={modalMoveOpen}
        handleCloseModal={() => setModalMoveOpen(false)}
        handleSubmitModal={() => setModalMoveOpen(false)}
      >
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Select source area</Form.Label>
            <Form.Select onChange={(e) => setMoveSource(e.target.value)}>
              <option>Please select area</option>
              <option value="1">Organic lettuce</option>
            </Form.Select>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Select destination area</Form.Label>
            <Form.Select onChange={(e) => setMoveDest(e.target.value)}>
              <option>Please select area</option>
              <option value="1">Organic lettuce</option>
              <option value="2">Organic chilli</option>
            </Form.Select>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>{`How many plants do you want to move? (${moveQty})`}</Form.Label>
            <Form.Range
              value={moveQty}
              onChange={(e) => setMoveQty(Number(e.target.value))}
            />
          </Form.Group>
        </Form>
      </ModalContainer>

      <ModalContainer
        title="Dump rom-24mar"
        isShow={modalDumpOpen}
        handleCloseModal={() => setModalDumpOpen(false)}
        handleSubmitModal={() => setModalDumpOpen(false)}
      >
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Choose area</Form.Label>
            <Form.Select onChange={(e) => setDumpArea(e.target.value)}>
              <option>Please select area</option>
              <option value="1">Organic lettuce</option>
            </Form.Select>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>{`How many plants do you want to dump? (${dumpQty})`}</Form.Label>
            <Form.Range
              value={moveQty}
              onChange={(e) => setDumpQty(Number(e.target.value))}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Notes</Form.Label>
            <Form.Control
              as="textarea"
              onChange={(e) => setDumpNotes(e.target.value)}
              value={dumpNotes}
              style={{ height: "120px" }}
            />
          </Form.Group>
        </Form>
      </ModalContainer>

      <ModalContainer
        title="Take Picture"
        isShow={modalPicOpen}
        handleCloseModal={() => setModalPicOpen(false)}
        handleSubmitModal={() => setModalPicOpen(false)}
      >
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Choose photo</Form.Label>
            <Form.Control
              type="file"
              onChange={(e) => setPic(e.target.value)}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>
              Photo description
              <br />
              <small className="text-muted">(max. 200 char)</small>
            </Form.Label>
            <Form.Control
              as="textarea"
              onChange={(e) => setPicNotes(e.target.value)}
              value={picNotes}
              style={{ height: "120px" }}
            />
          </Form.Group>
        </Form>
      </ModalContainer>
    </Layout>
  );
};

export default CropDetail;
