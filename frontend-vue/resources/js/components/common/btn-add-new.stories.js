import { storiesOf } from '@storybook/vue';
import BtnAddNew from './btn-add-new.vue';

storiesOf('BtnAddNew', module)
  .add('default reservoir', () => (
    {
      components: { BtnAddNew },
      template: '<BtnAddNew title="Add Reservoir" />',
    }
  ))
  .add('default area', () => (
    {
      components: { BtnAddNew },
      template: '<BtnAddNew title="Add Area" />',
    }
  ))
  .add('with custom css', () => (
    {
      components: { BtnAddNew },
      template: '<BtnAddNew title="Add Area" class="float-right" />',
    }
  ));
