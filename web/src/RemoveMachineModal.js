import React from 'react';
import { Button, Icon, Modal, Header } from 'semantic-ui-react'

export default ({isOpen, trigger, name, onRemove, onClose}) => (
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
