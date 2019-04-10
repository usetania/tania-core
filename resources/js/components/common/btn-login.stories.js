import { storiesOf } from '@storybook/vue';
import BtnLogin from './btn-login.vue';

storiesOf('BtnLogin', module)
  .add('default', () => (
    {
      components: { BtnLogin },
      template: '<BtnLogin />',
    }
  ));
