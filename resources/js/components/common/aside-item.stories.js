import { storiesOf } from '@storybook/vue';
import AsideItem from './aside-item.vue';

storiesOf('AsideItem', module)
  .add('default', () => (
    {
      components: { AsideItem },
      template: '<b-nav vertical><AsideItem title="Dashboard" fontawesome="fa fa-home" :routeName="{ name: \'Home\' }" /></b-nav>',
    }
  ))
  .add('active', () => (
    {
      components: { AsideItem },
      template: '<b-nav vertical><AsideItem title="Dashboard" fontawesome="fa fa-home" :routeName="{ name: \'Home\' }" :isActive=true /></b-nav>',
    }
  ));
