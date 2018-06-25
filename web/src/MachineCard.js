import React, { Component } from 'react'
import { Button, Card, Label } from 'semantic-ui-react'

import RemoveMachineModal from './RemoveMachineModal'

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

  export default MachineCard
