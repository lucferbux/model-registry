import React from 'react';
import {
  Dropdown,
  DropdownItem,
  DropdownList,
  Masthead,
  MastheadContent,
  MastheadMain,
  MenuToggle,
  MenuToggleElement,
  Toolbar,
  ToolbarContent,
  ToolbarGroup,
  ToolbarItem,
} from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { SimpleSelect, SimpleSelectOption } from '@patternfly/react-templates';
import { useDeepCompareMemoize } from '~/shared/utilities/useDeepCompareMemoize';

interface NavBarProps {
  username?: string;
  onLogout: () => void;
  options: SimpleSelectOption[];
  onNamespaceSelect?: (namespace: string) => void;
}

const NavBar: React.FC<NavBarProps> = ({ username, onLogout, options, onNamespaceSelect }) => {
  const [selected, setSelected] = React.useState<string | undefined>(String(options[0]?.value));
  const [userMenuOpen, setUserMenuOpen] = React.useState(false);
  const optionsMemo = useDeepCompareMemoize(options);
  const navigate = useNavigate();

  const initialOptions = React.useMemo<SimpleSelectOption[]>(
    () => optionsMemo.map((o) => ({ ...o, selected: o.value === selected })),
    [selected, optionsMemo],
  );

  React.useEffect(() => {
    if (selected) {
      navigate(`?namespace=${selected}`);
    }
  }, [selected, navigate]);

  const handleLogout = () => {
    setUserMenuOpen(false);
    onLogout();
  };

  const userMenuItems = [
    <DropdownItem key="logout" onClick={handleLogout}>
      Log out
    </DropdownItem>,
  ];

  return (
    <Masthead>
      <MastheadMain />
      <MastheadContent>
        <Toolbar>
          <ToolbarContent>
            <ToolbarGroup variant="action-group-plain" align={{ default: 'alignStart' }}>
              <ToolbarItem>
                <SimpleSelect
                  initialOptions={initialOptions}
                  onSelect={(_ev, selection) => {
                    setSelected(String(selection));
                    if (onNamespaceSelect) {
                      onNamespaceSelect(String(selection));
                    }
                  }}
                />
              </ToolbarItem>
            </ToolbarGroup>
            {username && (
              <ToolbarGroup variant="action-group-plain" align={{ default: 'alignEnd' }}>
                <ToolbarItem>
                  <Dropdown
                    popperProps={{ position: 'right' }}
                    onOpenChange={(isOpen) => setUserMenuOpen(isOpen)}
                    toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
                      <MenuToggle
                        aria-label="User menu"
                        id="user-menu-toggle"
                        ref={toggleRef}
                        onClick={() => setUserMenuOpen(!userMenuOpen)}
                        isExpanded={userMenuOpen}
                      >
                        {username}
                      </MenuToggle>
                    )}
                    isOpen={userMenuOpen}
                  >
                    <DropdownList>{userMenuItems}</DropdownList>
                  </Dropdown>
                </ToolbarItem>
              </ToolbarGroup>
            )}
          </ToolbarContent>
        </Toolbar>
      </MastheadContent>
    </Masthead>
  );
};

export default NavBar;
