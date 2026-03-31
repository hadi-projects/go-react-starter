import { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import PropTypes from "prop-types";

const Sidebar = ({
  sections = [],
  title = "Admin Panel",
  onLogout,
  collapsed = false,
  onToggleCollapse,
}) => {
  const location = useLocation();
  const [expandedSections, setExpandedSections] = useState({});

  const isActive = (path) => location.pathname === path;
  const isChildActive = (item) => {
    if (!item.subItems) return isActive(item.path);
    return item.subItems.some((sub) => isActive(sub.path));
  };

  // Auto-expand section if a child is active (only when not collapsed)
  useEffect(() => {
    if (collapsed) return;
    const newExpanded = { ...expandedSections };
    let changed = false;
    sections.forEach((section) => {
      section.items.forEach((item) => {
        if (
          item.subItems &&
          isChildActive(item) &&
          !expandedSections[item.label]
        ) {
          newExpanded[item.label] = true;
          changed = true;
        }
      });
    });
    if (changed) setExpandedSections(newExpanded);
  }, [location.pathname, sections, collapsed]);

  const toggleSection = (label) => {
    if (collapsed) return;
    setExpandedSections((prev) => ({ ...prev, [label]: !prev[label] }));
  };

  return (
    <aside
      className={`nav-drawer flex-shrink-0 flex flex-col transition-all duration-300 ease-in-out ${collapsed ? "w-[56px]" : "w-56"}`}
      style={{ overflow: "hidden" }}
    >
      {/* Header with toggle button */}
      <div
        className={`h-12 flex items-center flex-shrink-0 ${collapsed ? "justify-center px-2" : "justify-between px-3"}`}
      >
        {!collapsed && (
          <h2 className="text-sm font-semibold text-surface-on truncate">
            {title}
          </h2>
        )}
        <button
          onClick={onToggleCollapse}
          className="p-1.5 rounded-full hover:bg-surface-variant/40 text-surface-on-variant transition-all duration-200 flex-shrink-0"
          title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          <svg
            className={`w-4 h-4 transition-transform duration-300 ${collapsed ? "rotate-180" : ""}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M11 19l-7-7 7-7m8 14l-7-7 7-7"
            />
          </svg>
        </button>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto overflow-x-hidden py-2 custom-scrollbar">
        {sections.map((section, idx) => (
          <div key={idx} className="mb-4">
            {section.label && !collapsed && (
              <h3 className="px-4 py-1.5 text-[10px] font-semibold text-surface-on-variant uppercase tracking-wider whitespace-nowrap">
                {section.label}
              </h3>
            )}
            {section.label && collapsed && (
              <div className="px-2 py-2">
                <div className="border-t border-outline-variant/30" />
              </div>
            )}
            <ul className={`space-y-0.5 ${collapsed ? "px-1.5" : "px-2"}`}>
              {section.items.map((item) => {
                const active = isChildActive(item);
                const hasSubItems = !!item.subItems;
                const isExpanded = expandedSections[item.label];

                if (collapsed) {
                  // Icon-only mode: flat, no sub-items expansion
                  const href = item.subItems
                    ? item.subItems[0]?.path
                    : item.path;
                  return (
                    <li key={item.label}>
                      <Link
                        to={href || "#"}
                        title={item.label}
                        className={`flex items-center justify-center w-full p-2 rounded-full transition-all duration-200 ${
                          active
                            ? "bg-secondary-container text-secondary-on-container"
                            : "text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on"
                        }`}
                      >
                        {item.icon}
                      </Link>
                    </li>
                  );
                }

                return (
                  <li key={item.label}>
                    {hasSubItems ? (
                      <div>
                        <button
                          onClick={() => toggleSection(item.label)}
                          className={`w-full flex items-center justify-between px-3 py-2 rounded-full transition-all duration-300 group ${
                            active
                              ? "bg-secondary-container/50 text-secondary-on-container"
                              : "text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on"
                          }`}
                        >
                          <div className="flex items-center gap-2">
                            <div className="transition-transform duration-300 group-hover:scale-110 flex-shrink-0">
                              {item.icon}
                            </div>
                            <span className="font-medium text-xs whitespace-nowrap">
                              {item.label}
                            </span>
                          </div>
                          <svg
                            className={`w-3.5 h-3.5 flex-shrink-0 transition-transform duration-300 ${isExpanded ? "rotate-180" : ""}`}
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M19 9l-7 7-7-7"
                            />
                          </svg>
                        </button>

                        <div
                          className={`overflow-hidden transition-all duration-300 ease-in-out ${isExpanded ? "max-h-64 mt-0.5" : "max-h-0"}`}
                        >
                          <ul className="space-y-0.5 ml-3 pl-3 border-l border-outline-variant/30">
                            {item.subItems.map((sub) => {
                              const subActive = isActive(sub.path);
                              return (
                                <li key={sub.path}>
                                  <Link
                                    to={sub.path}
                                    className={`flex items-center gap-2 px-3 py-1.5 rounded-full transition-all duration-300 group ${
                                      subActive
                                        ? "bg-secondary-container text-secondary-on-container shadow-sm"
                                        : "text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on"
                                    }`}
                                  >
                                    <div className="transition-transform duration-300 group-hover:scale-110 flex-shrink-0">
                                      {sub.icon}
                                    </div>
                                    <span className="font-medium text-xs whitespace-nowrap">
                                      {sub.label}
                                    </span>
                                  </Link>
                                </li>
                              );
                            })}
                          </ul>
                        </div>
                      </div>
                    ) : (
                      <Link
                        to={item.path}
                        className={`flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-300 group ${
                          active
                            ? "bg-secondary-container text-secondary-on-container"
                            : "text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on"
                        }`}
                      >
                        <div className="transition-transform duration-300 group-hover:scale-110 flex-shrink-0">
                          {item.icon}
                        </div>
                        <span className="font-medium text-xs whitespace-nowrap">
                          {item.label}
                        </span>
                      </Link>
                    )}
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>

      {/* Bottom: Logout */}
      <div
        className={`flex-shrink-0 p-1.5 border-t border-outline-variant/30 ${collapsed ? "flex justify-center" : "px-2"}`}
      >
        <button
          onClick={onLogout}
          title="Logout"
          className={`flex items-center gap-2 px-3 py-2 rounded-full text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on transition-all duration-300 group ${collapsed ? "justify-center w-full" : "w-full"}`}
        >
          <div className="transition-transform duration-300 group-hover:scale-110 flex-shrink-0">
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
              />
            </svg>
          </div>
          {!collapsed && (
            <span className="font-medium text-xs whitespace-nowrap">
              Logout
            </span>
          )}
        </button>
      </div>
    </aside>
  );
};

Sidebar.propTypes = {
  sections: PropTypes.arrayOf(
    PropTypes.shape({
      label: PropTypes.string,
      items: PropTypes.arrayOf(
        PropTypes.shape({
          path: PropTypes.string,
          label: PropTypes.string.isRequired,
          icon: PropTypes.node.isRequired,
          subItems: PropTypes.arrayOf(
            PropTypes.shape({
              path: PropTypes.string.isRequired,
              label: PropTypes.string.isRequired,
              icon: PropTypes.node.isRequired,
            }),
          ),
        }),
      ).isRequired,
    }),
  ),
  title: PropTypes.string,
  onLogout: PropTypes.func.isRequired,
  collapsed: PropTypes.bool,
  onToggleCollapse: PropTypes.func.isRequired,
};

export default Sidebar;
