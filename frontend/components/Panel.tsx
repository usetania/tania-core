import { Card, Col } from "react-bootstrap";

interface iPanel {
  title: string;
  children: JSX.Element;
  xs?: number;
  sm?: number;
  md?: number;
  lg?: number;
  footer?: string;
}

const Panel = ({
  title,
  children,
  xs = 12,
  sm = 12,
  md = 12,
  lg = 12,
  footer,
}: iPanel): JSX.Element => {
  return (
    <Col xs={xs} sm={sm} md={md} lg={lg}>
      <Card className="mb-3">
        <Card.Body>
          <Card.Title>
            <h5>{title}</h5>
          </Card.Title>
          {children}
        </Card.Body>
        {footer && <Card.Footer>{footer}</Card.Footer>}
      </Card>
    </Col>
  );
};

export default Panel;
