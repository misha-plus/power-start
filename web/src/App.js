import React, { Component } from 'react';
import './App.css';
import 'semantic-ui-css/semantic.min.css';
import {
  Button, Card, Image, Dimmer, Loader, Icon, Modal, Header, Label
} from 'semantic-ui-react'
import axios from 'axios';

const RemoveMachineModal = ({isOpen, trigger, name, onRemove, onClose}) => (
  <Modal open={isOpen} basic size='small' onClose={onClose}>
    <Header icon='trash alternate' content={`Removing machine "${name}"`} />
    <Modal.Content>
      <p>
        Do you want to remove machine "{name}" from power-start?
        You can't undo this action. Machine will keep current power state.
      </p>
    </Modal.Content>
    <Modal.Actions>
      <Button basic color='red' inverted onClick={onRemove}>
        <Icon name='trash alternate' /> Remove
      </Button>
      <Button color='green' inverted onClick={onClose}>
        <Icon name='checkmark' /> Cancel
      </Button>
    </Modal.Actions>
  </Modal>
)

class MachineCard extends Component {
  constructor ({ name, mac, requests, isRunning, onStart, onStop, onRemove }) {
    super()
    this.props = { name, mac, requests, isRunning, onStart, onStop, onRemove }
    this.state = { isRemoveModalOpen: false }
  }

  render() {
    const [, color, status] = [
      [(x, y) => x == false && y == 0, 'red',    'Stopped'],
      [(x, y) => x == false && y > 0,  'orange', 'Pending run'],
      [(x, y) => x == true  && y > 0,  'green',  'Running'],
      [(x, y) => x == true  && y == 0, 'orange', 'Pending stop'],
    ].find(row => row[0](this.props.isRunning, this.props.requests))

    return (
      <Card color={color}>
        <Card.Content>
          <Button
            floated='right'
            circular
            icon='trash alternate'
            onClick={() => this.setState({isRemoveModalOpen: true})}
          />
          <RemoveMachineModal
            isOpen={this.state.isRemoveModalOpen}
            onClose={() => this.setState({isRemoveModalOpen: false})}
            onRemove={() => {
              this.setState({isRemoveModalOpen: false})
              this.props.onRemove()
            }}
            name={this.props.name}
          />
          <Card.Header>{this.props.name}</Card.Header>
          <Card.Meta>MAC: {this.props.mac}</Card.Meta>
          <Card.Description>
            <p>Requests: {this.props.requests}</p>
            <p>
              <Label color={color} horizontal>{status}</Label>
            </p>
          </Card.Description>
        </Card.Content>
        <Card.Content extra className="running-machine-card">
          <div className='ui two buttons'>
            <Button basic color='green' onClick={this.props.onStart}>
              Start
            </Button>
            <Button
              basic
              color='red'
              onClick={this.props.onStop}
              disabled={this.props.requests == 0}
            >
              Stop
            </Button>
          </div>
        </Card.Content>
      </Card>
    )
  }
}


class App extends Component {
  constructor() {
    super()
    this.state = { data: null }
    this.updateData()
    setInterval(this.updateData.bind(this), 2000)
  }

  async updateData() {
    const resp = await axios.get('/api/list')
    if (resp.status != 200) {
      alert("Error while request new data") // TODO
      return
    }
    this.setState({ data: resp.data })
  }

  async start(name) {
    await axios.post(`/api/start/${name}`)
    await this.updateData()
  }

  async stop(name) {
    await axios.post(`/api/stop/${name}`)
    await this.updateData()
  }

  async remove(name) {
    alert('removing ' + name)
    await axios.post(`/api/remove/${name}`)
    await this.updateData()
  }

  render() {
    if (!this.state.data) {
      return (
        <Dimmer active>
          <Loader>
            Loading machines list...
            <p>Power-Start</p>
          </Loader>
        </Dimmer>
      )
    }

    return (
      <div>
        <div className='power-start-header'>
          <Header as='h2' icon textAlign='center'>
            <Image src='icon.png' size='massive' circular />
            <Header.Content>Power-Start</Header.Content>
          </Header>
        </div>
        <Card.Group className='power-start-cards'>
          {this.state.data.map(machine =>
            <MachineCard
              key={machine.name}
              name={machine.name}
              mac={machine.mac}
              requests={machine.requests}
              isRunning={machine.isRunning}
              onStart={(n => () => this.start(n))(machine.name)}
              onStop={(n => () => this.stop(n))(machine.name)}
              onRemove={(n => () => this.remove(n))(machine.name)}
            />
          )}
        </Card.Group>
      </div>
    );
  }
}

export default App;
