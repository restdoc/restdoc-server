// nav menu 
    document.addEventListener('DOMContentLoaded', () => {

    // Get all "navbar-burger" elements
    const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);

    // Check if there are any navbar burgers
    if ($navbarBurgers.length > 0) {

    // Add a click event on each of them
    $navbarBurgers.forEach( el => {
      el.addEventListener('click', () => {

        // Get the target from the "data-target" attribute
        const target = el.dataset.target;
        const $target = document.getElementById(target);

        // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
        el.classList.toggle('is-active');
        $target.classList.toggle('is-active');

      });
    });
  }

});

// track
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'UA-143200156-1');

  //current menu add css
  var path = window.location.pathname;
  path = path.trim();
  var eid = "";
  var menuId = path.substr(1, path.length) 
  if (menuId == "") {
	  menuId = "home";
  }

  if (menuId){
	var arr = menuId.split('/');
	var hasLevel = false;
	if (arr.length > 1) {
		hasLevel = true;
		menuId = arr[0];
	}


	var menuName = menuId + "Menu";
	var menu = document.getElementById(menuName);
	if (menu && menu != null && menu != undefined){
		if (hasLevel) {
			menu.classList.add("current-sub");
		} else {
			menu.classList.add("current");
		}
	}
  }
