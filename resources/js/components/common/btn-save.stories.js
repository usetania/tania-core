import { storiesOf } from '@storybook/vue';
import BtnSave from './btn-save.vue';

storiesOf('BtnSave', module)
  .add('default', () => (
    {
      components: { BtnSave },
      template: '<BtnSave />',
    }
  ))
  .add('float-right', () => (
    {
      components: { BtnSave },
      template: '<BtnSave customClass="float-right" />',
    }
  ));
